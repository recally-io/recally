# 个人阅读档案：Cloudflare 架构与实施 Plan

**以 Pi 为 Agent runtime：`Skill 提供策略，Tools 提供能力，Agent 自主决定获取路径`。用户提交 URL 后，CaptureAgent 可按证据选择 HTTP fetch、站点适配器、Browser Run 或进一步读取，再提交可信快照；随后由 AnalysisAgent 完成总结与结构化理解。默认全 Cloudflare，不引入 Go 服务或自建浏览器。**

版本：v1.1  
日期：2026-09-11  
状态：供实现与任务拆分使用的设计方案，尚未实现或完成线上验证。  
产品名称：暂用“个人阅读档案”（候选品牌 ReKeep 尚未定稿）。

本文中的架构、预算、数据结构和验收门槛是**设计决定或初始配置**。Cloudflare 能力与限制以文中官方来源为依据，查阅日期为 2026-09-11。模型质量、抓取成功率、延迟和费用均需通过 M0 与真实样本验证，不作为已达到的指标。

## 阅读导航

**产品与架构**：第 1–4 节。  
**Pi Agent、Skills、获取工具、归档与可靠执行**：第 5–9 节。  
**AI 加工、搜索、界面与 API**：第 10–13 节。  
**安全、部署、预算与验证**：第 14–18 节。  
**里程碑与实施任务**：第 19 节。  
**发布标准、设计决策与来源**：第 20–22 节。

---

## 1. 产品目标与边界

### 1.1 要解决的问题

每天阅读的文章、项目介绍与网页链接分散在浏览器和 AI 对话中。需要一个统一入口，保存看过的内容、当时得到的 AI 总结与自己的笔记；以后即使原链接失效，也能阅读、搜索、回顾和有选择地分享。

产品的核心对象是**个人阅读档案**，不是只有 URL 的书签列表，也不是面向整个互联网的搜索引擎。

### 1.2 本次明确采用的原则

**Agent 从第一版就是控制器，不是固定 pipeline 上的一次 LLM 判断。** 用户给出 URL 后，代码只创建 CaptureRun、加载 skill、暴露工具和预算；具体使用哪个获取策略、是否继续、是否切换策略、何时认为证据足够，由 Pi CaptureAgent 根据 observation 决定。

**Skills encode strategy; tools expose capability; agent owns orchestration。** `capture/SKILL.md` 描述“什么情况下适合 document fetch、site adapter、browser、何时停止”等经验；`web_fetch / site_run / browser_* / archive_*` 只执行单一能力。工具内部不得偷偷把 `fetch → browser → Jina` 串成固定 fallback，否则 Agent 实际看不到获取决策。[S31]

**`fetch first` 是默认偏好，不是代码硬约束。** 对 document-like 页面，capture skill 通常指导 Agent 先用低成本 HTTP；对 X/Reddit/GitHub 这类 app-like record，如果已有可靠 site adapter，Agent 可以直接使用 adapter；对用户已提交正文或明显需要渲染的页面，也可以跳过无意义的 fetch。策略变化通过 skill revision 发布，不靠不断增加 host-specific `if/else`。

**Pi 是可移植的 Agent runtime。** MVP 直接使用 `@earendil-works/pi-agent-core`；不使用完整 `pi-coding-agent` 的 bash/filesystem/coding 工具链。Pi core 提供 tool calling 与 state management，本项目自行定义 ReKeep 的 skills、tools、state store 和权限边界。[S32]

**Capture 与 Analysis 是两个职责不同的 Agent。** CaptureAgent 的唯一目标是“忠实获取并提交用户想保存的内容”；AnalysisAgent 只读取已归档 Snapshot，负责总结、关键观点、标签与引用。分析失败不能反过来影响原文归档。

**AI 总结与原文分开保存。** AI 可以选择正文、提出保存方案、生成摘要和标签；原始响应、rendered DOM、结构化站点记录、选中的原文块与历史版本不能被模型自由改写。

**Workflow 负责 durable execution，不负责业务策略。** Workflow 保存 Agent turn、工具结果、重试和恢复；它不写“先 fetch、失败 browser、然后 summary”的业务分支。决策属于 Pi，恢复与副作用语义属于 Workflow。

**默认全 Cloudflare。** Workers、Workflows、D1、R2、Browser Run、Workers AI、Vectorize 构成完整生产路径。平台相关部分通过 adapters 隔离，使 CaptureAgent/skills 未来可以复用于 celld self-hosted 版本。

**个人使用起步，默认私有。** 数据包含 `library_id`，但第一版不做团队协作、开放注册或多租户运营平台。

### 1.3 MVP 的交付范围

**采集与保存**：提交 URL；AI 驱动的 fetch/browser 获取；保存正文、来源和关键资源；失败时可粘贴正文补充；支持手动重新归档与版本查看。

**理解与管理**：AI 摘要、关键观点、主题标签；手动粘贴已有 AI 总结；个人笔记、收藏夹、阅读状态与时间线。

**找回与使用**：中英文全文搜索、基本语义搜索、带原文引用的知识库问答、选定条目的可撤销分享、可下载导出。

**可靠性**：重试、预算、可取消任务、删除不复活、权限检查、备份恢复和基本可观测性。

### 1.4 后续阶段，不能成为 MVP 依赖

浏览器扩展、原生手机应用、自动浏览历史同步、复杂划线与协作、第三方阅读服务批量导入、RSS、知识图谱、推荐系统、站点级爬虫、视频下载与转录、OCR、完整 WARC 重放均后置。

第一版可以保存遇到的非 HTML 文件为附件，但不能声称已完成其正文解析、AI 总结或全文检索。PDF、X 长线程等专用适配器另设后续里程碑。

### 1.5 用户层面的成功标准

用户保存一篇文章后，能够确认保存的是正文而不是登录页；能分辨“完整”“部分”“待处理”；源站失效后能读已保存内容；能用模糊描述找回文章；AI 回答可以跳回原文；分享不会夹带私人笔记。

**“永久保存”是产品方向，不是不可失效的承诺。** 实际可用性依赖成功抓取、存储续费、备份、权限与恢复能力。系统必须允许完整导出，避免把数据可读性绑定到运行平台。

## 2. 关键用户流程

### 2.1 Document-like 页面：Agent 通常选择 fetch

用户提交普通博客或文档 URL → 创建 CaptureRun → 加载 capture skill → Pi 查看 URL 与可用 tools → 通常调用 `web_fetch` → 工具执行一次 HTTP 请求与确定性正文提取 → 返回标题、正文块、长度、结构、重定向、截断等 observation → Agent 判断内容是否代表用户想保存的页面。

如果证据充分，Agent 直接 `propose_archive`；如果正文可疑，它可以继续 `read_source`、尝试另一确定性 extractor、寻找观察到的打印版/下一页，或升级 browser。**“通常先 fetch”来自 skill 的成本与可靠性经验，不来自业务代码的强制分支。**

### 2.2 App-like URL：Agent 可以优先选择 site adapter

对于一条 Tweet、Reddit post、GitHub repo、HN item、视频详情等，用户想保存的往往是一个结构化 record，而不是页面 chrome。CaptureAgent 先通过 `site_list(url)` 查看是否存在适配器；若 skill 判断 adapter 更忠实，可以直接 `site_run(adapter, args)`，得到结构化 JSON、媒体引用和 canonical URL。

```text
URL
 ↓
CaptureAgent + capture skill
 ↓
site_list
 ↓
site_run("twitter/status", ...)
 ↓
structured record + media refs
 ↓
Agent inspect
 ↓
archive
```

这继承 Stella web skill 的核心经验：document 用 fetch，application record 优先站点自己的数据路径；以后增加来源主要新增 adapter/skill，而不是改 Agent 主循环。[S31]

### 2.3 动态页面：Agent 自主进入 Browser Run

HTTP/adapter 结果为空、只含 JS shell、正文未展开，或 Agent 根据 evidence 判断必须交互时，调用 browser capability。浏览器逻辑仍由 Agent 决定：observe → act → observe，可以等待正文、滚动、展开、导航同一篇文章的下一页，再决定 capture 或退出。

Browser Run 只是工具执行环境，不是另一个硬编码 scraper。逻辑上 Pi 看到的是 `browser_open / browser_observe / browser_act / browser_capture`；实现上为了保持 page handle，这些动作可以在一个 Workflow browser episode 内连续执行，但每个动作仍由 Pi 根据最新 observation 决定。[S02][S03]

### 2.4 登录墙、验证码或不可访问内容

Agent 与确定性规则发现访问限制 → 停止自动交互 → 保留收藏、已取得证据和失败原因 → 提示用户粘贴有权保存的正文，或后续由浏览器扩展主动提交当前页面。

不得用“Agent 更灵活”作为绕过付费、登录、验证码或访问权限的理由。Browser Run 也不保证通过站点反爬。[S04]

### 2.5 Capture 完成后由 AnalysisAgent 接手

Snapshot 提交成功后，AnalysisAgent 加载 `analyze` skill，只读取归档内容和稳定 block ID。它可以按需读取正文范围、生成局部摘要、检查遗漏，再提交 summary/key points/topics/entities/citations。它没有 browser、任意 web fetch、权限修改或删除工具，因此总结过程不会悄悄把实时网页混入历史快照。

### 2.6 原文已归档，但模型调用失败

用户仍可阅读与全文搜索；界面显示“总结待重试”或“语义索引待重试”。重试从相应加工阶段开始，不重新抓取已经确认的快照。

### 2.7 再次收藏同一 URL

普通重复提交返回已有 Item；显式“重新抓取”创建新的 CaptureRun。新的 Agent run 可以采用与旧 run 不同的 skill/tool/model revision；新内容产生新 Snapshot，内容相同可以复用对象，但仍保留本次抓取时间、Agent decisions 与来源证据。

旧总结、旧引用和已发布分享继续指向原版本，不随最新网页改变。

## 3. 总体架构

```text
Web / URL API / Browser Extension（后续）
                    │
                    ▼
┌──────────────────── Cloudflare ──────────────────────┐
│                                                     │
│ App Worker + React Static Assets                    │
│  认证 / 收藏 / 阅读 / 笔记 / 搜索 / 分享              │
│                    │                                │
│              create CaptureRun                      │
│                    ▼                                │
│              IngestWorkflow                         │
│        durable supervisor / checkpoints             │
│                    │                                │
│                    ▼                                │
│          ┌─────────────────────┐                    │
│          │   Pi CaptureAgent   │                    │
│          │ @pi-agent-core      │                    │
│          │ + capture/SKILL.md  │                    │
│          └──────────┬──────────┘                    │
│                     │ decides                       │
│       ┌─────────────┼──────────────┐                │
│       ▼             ▼              ▼                │
│   web_fetch      site_run      browser_*            │
│ fetch+Defuddle   adapters      Browser Run           │
│       │             │              │                │
│       └─────────────┼──────────────┘                │
│                     ▼                               │
│              observations / evidence                │
│                     │                               │
│                  Pi decides                         │
│                     │                               │
│             propose_archive                         │
│                     ▼                               │
│            ArchiveService ───────────────► R2       │
│                     │                               │
│                     ├──────────────► D1             │
│                     │    product state / FTS        │
│                     ▼                               │
│              EnrichWorkflow                         │
│                     │                               │
│          ┌─────────────────────┐                    │
│          │  Pi AnalysisAgent   │                    │
│          │ + analyze/SKILL.md  │                    │
│          └──────────┬──────────┘                    │
│                     ├────────────► Workers AI       │
│                     └────────────► Vectorize        │
│                                                     │
│ Cron：outbox 补发、回顾、清理、对账、备份             │
└─────────────────────────────────────────────────────┘
```

### 3.1 Agent runtime

MVP 的 Agent runtime 使用 `@earendil-works/pi-agent-core`：它提供 stateful agent、tool execution 和 event streaming，并将 Node SQLite backend 从 core 拆开，因此可以为 Worker runtime 提供自己的 state adapter。[S32]

**不采用完整 coding agent。** ReKeep 不需要 bash、git、任意文件系统或代码执行。CaptureAgent 的能力只能来自显式注册的 Web/Archive tools。

### 3.2 部署单元

**App Worker**：用户 API、私有阅读页面、搜索与问答、分享授权；使用 Workers Static Assets 托管 React 构建结果。[S01]

**Jobs Worker**：Pi runtime、skills、tool registry、Workflow 类、Browser Run adapter、计划任务。Browser binding 只配置给 Jobs Worker，避免普通搜索与页面请求具备浏览器能力。

两者位于同一 TypeScript monorepo。第一版不需要 Go/Python 服务、自建 Chromium、Redis、NATS 或 Temporal。

### 3.3 Cloudflare 产品的职责

**D1** 保存产品业务状态、Agent run 索引、授权与 FTS；**R2** 保存不可变网页证据、Agent 大型 observations、快照、附件与导出；**Vectorize** 是可重建的语义索引。

**Workflows** 提供 durable supervisor、step retry 和恢复，不决定获取策略；**Workers AI / Pi model provider** 提供 Capture/Analysis/Answer 模型；**Browser Run** 执行完整浏览器工具。Cloudflare Browser Run 支持 Worker 内的 Playwright 与直接浏览器会话。[S02][S03]

### 3.4 Skill 与 Tool 的边界

`skills/capture/SKILL.md` 是版本化的策略知识：告诉 Agent 如何区分 document、application record、dynamic page、受限内容，如何评估“获取是否足够”，以及成本/安全原则。它不包含站点凭据，也不直接执行网络动作。

Tool registry 提供版本化能力：`web_fetch`、`site_list/site_run`、`browser_*`、`read_source`、`archive_*`。Tool 的返回全部视为不可信 observation。**新增来源优先通过 adapter + skill 扩展，而不是改 CaptureAgent 主循环。**

### 3.5 平台可移植边界

Agent core、skills、tool contracts、domain model 不直接引用 Cloudflare binding。Cloudflare 实现放在 platform adapters：

```text
CaptureAgent / Skills / Contracts
             │
      Platform Adapters
       /            \
 Cloudflare        future celld
 Browser Run       Playwright service
 R2                S3/R2
 D1                Cell SQLite / DB
 Vectorize         sqlite-vec / external
 Workflows         celld Workflow
```

MVP 只实现 Cloudflare adapter，不为了未来 self-hosted 提前复制两套系统；但避免让 skill 或 Agent loop 直接依赖 Browser Run、R2 key 规则或 D1 SQL。

### 3.6 第一版不引入的组件

不引入 Cloudflare Agents SDK 作为必选 Agent runtime，也不引入 Queue、KV 或通用插件市场。Pi 是 Agent runtime，Workflow/D1/R2 是执行和存储底座。

不同时维护 AI Search 与自建检索链。本版采用 D1 FTS5 + Vectorize，以直接控制快照版本、正文块和引用关系。

## 4. 领域模型与存储设计

### 4.1 核心实体

```text
Library
 ├── Item
 │    ├── CaptureRun
 │    │    ├── CaptureAttempt
 │    │    └── AgentEvent / Observation
 │    ├── Snapshot
 │    │    ├── SourceDocument / Asset
 │    │    └── ContentRevision
 │    │         └── Chunk
 │    ├── AIArtifact
 │    ├── Note / ReadingEvent
 │    └── Collection / Tag memberships
 ├── Share
 ├── Job / Outbox
 ├── API Token
 └── UsageReservation / UsageEvent
```

**Item**：一条收藏的稳定身份。主要字段为 `id, library_id, original_url, normalized_url, title, saved_at, first_read_at, last_read_at, read_status, current_snapshot_id, capture_generation, deleted_at`。

**CaptureRun / CaptureAttempt**：一次用户要求的抓取及其具体执行尝试。记录 `job_id, agent_runtime_version, skill_revision, toolset_version, policy_version, pipeline_version, model_config_version, requested_at, target_url, status, outcome_code, budget, generation`。同一个 URL 在不同 skill revision 下可以走不同策略；网络重试不能伪装成用户重新收藏。

**Snapshot**：某次已提交的来源快照。记录 `captured_at, final_url, capture_method, source_manifest_key, manifest_hash, content_quality, resource_quality, committed_at`。一个跨页文章包含多个 SourceDocument，每页分别保存来源与时间。

**ContentRevision**：从已保存 Snapshot 提取出的阅读正文版本。记录 `extractor_version, selected_source_blocks, body_hash, article_key, blocks_key, language, token_count`。改进正文提取不必重新访问互联网，也不能覆盖旧锚点。

**Chunk**：用于检索的分块。记录 `content_revision_id, ordinal, block_range, text, text_hash, tokenizer_version, embedding_version, vector_id, index_state`。

**AIArtifact**：总结、关键观点、导入的 AI 回答或回顾。记录输入 Snapshot/ContentRevision、模型、提示词版本、覆盖范围、输出、引用、费用和生成时间。

**Note / ReadingEvent**：用户自己写的内容和实际阅读事件。笔记带版本，不能被自动总结覆盖；收藏不自动等于阅读，时区保存在 Library。

### 4.2 D1 表与主要约束

建议建表：`libraries, items, item_urls, capture_runs, capture_attempts, snapshots, source_documents, content_revisions, assets, chunks, ai_artifacts, notes, note_revisions, collections, item_collections, tags, item_tags, reading_events, shares, jobs, outbox, idempotency_records, agent_events, api_tokens, usage_reservations, usage_events, resource_slots, library_maintenance, deletion_tombstones, backup_runs`。

这是表结构蓝图，实际 migration 在 M1 实现；不是要求每张表对应一个服务。

**归属约束**：所有业务对象带 `library_id`；关联写入验证同一归属。关键外键采用能够验证 `(library_id, object_id)` 的组合约束，不能只凭跨库可猜测的 ID 关联。

**去重约束**：对未删除条目的有效 URL 映射，`item_urls(library_id, normalized_url_hash)` 唯一；URL 原文保留。删除事务同步停用旧映射，用户再次收藏可创建新 Item ID，不复活旧任务或分享。HTTP `Idempotency-Key` 作用域包含 actor、route 和 library，并记录 payload hash；同一 key 不同输入返回 409。

**版本约束**：Snapshot 和 ContentRevision 提交后不可变。相同输入与处理版本的 AIArtifact/Index job 通过 operation key 去重，不按“每个 Item 只有一个总结”覆盖。

**任务约束**：`operation_key`、`workflow_instance_id` 唯一；AgentEvent 的 `(run_id, attempt_id, sequence)` 唯一。发布结果时检查 Item 未删除且 generation 未过期。

**检索约束**：FTS 索引正文块与条目字段，返回前联查有效归属、删除状态和内容版本。D1 支持 FTS5。[S09]

**索引建议**：`items(library_id, saved_at, id)`、`jobs(status, next_attempt_at, id)`、`snapshots(library_id, item_id, captured_at)`、`chunks(library_id, content_revision_id, ordinal)`、`shares(token_hash)`、`reading_events(library_id, occurred_at)`。

### 4.3 URL 规范化

保留 original URL 与最终跳转 URL。规范化只处理 scheme/host 大小写、默认端口和明确的追踪参数，不删除内容相关参数，也不把任意 `canonical` 标签当成授权依据。

普通页面锚点可作为阅读位置保存；可能表示 SPA 路由的 fragment 必须保留语义。跨域 canonical 或跳转要重新经过网络策略，并在合并 Item 前检查来源一致性。

### 4.4 R2 组织

```text
archive/
  libraries/{library_id}/
    runs/{run_id}/attempts/{attempt_id}/
      observations/{observation_id}.json
      sources/{source_id}/response-body.html
      sources/{source_id}/response-metadata.json
      sources/{source_id}/rendered.html
      sources/{source_id}/screenshot.webp
      decisions/{sequence}.json
    snapshots/{snapshot_id}/
      manifest.json
      resources/{asset_id}
    content/{content_revision_id}/
      blocks.jsonl
      article.md
      provenance.json
    artifacts/{artifact_id}.json
    exports/{export_id}/...
```

已接受快照可以引用 run 下的不可变 source 对象，避免重复搬运；manifest 明确列出全部对象。**GC 不能简单按 `runs/` 路径过期删除**，必须先检查对象是否仍被已提交 Snapshot、导出或保留任务引用。

未被引用的失败尝试默认保留 7 天；已提交快照默认保留到用户删除。正文提取版本与 AI 结果按版本分别保存。

### 4.5 哪些数据可重建

全文索引、向量索引、模型生成标签可以从保存的输入重新构建。用户笔记、已读记录、历史分享配置和用户已经获得的 AI 总结是产品数据，**即使理论上可以重新调用模型，也应备份原版本**。

## 5. Pi CaptureAgent：Skill 驱动的内容获取

### 5.1 核心原则

CaptureAgent 的目标只有一句话：**忠实、可追溯、在预算内保存用户这个 URL 所指向的内容。**

它不是 `fetch()` 的包装器，也不是“固定 pipeline 中让 LLM 判断一个 boolean”。Pi 每一轮都会看到当前目标、capture skill、可用 tools、已取得 observations、预算和安全约束，然后决定下一步调用哪个 tool，或结束任务。

```text
Goal + Skill + Current State
            │
            ▼
      Pi CaptureAgent
            │
         decide
            │
   ┌────────┼─────────┐
   ▼        ▼         ▼
 web_fetch site_run browser_*
   │        │         │
   └────────┼─────────┘
            ▼
       observation
            │
            └────────► Pi next turn
                         │
                      ...
                         │
                 propose_archive / finish
```

### 5.2 为什么采用 Pi

`@earendil-works/pi-agent-core` 提供通用 Agent runtime、tool execution、state management 与事件流，本身不要求使用 coding-agent 的 bash/read/write/edit。它适合把本项目的受控 tools 注入进去。[S32]

**Plan 直接依赖 Pi Agent Core 的 agent/tool contract，并通过 ReKeep 自己的 adapters 适配 Cloudflare Workers。**

### 5.3 Capture Skill

Skill 采用 Stella/Agent Skills 的思路：Markdown 是 Agent 可读的“操作经验”，工具实现与策略文本分离。Stella 的 web skill 已明确区分 document fetch、site records 与 JS render，并把站点脚本作为可扩展能力；ReKeep 复用这一设计原则，但移除其 `web fetch` 内部固定的 HTTP → Lightpanda → Jina fallback。[S31]

建议结构：

```text
skills/
  capture/
    SKILL.md
    references/
      quality.md
      source-types.md
  analyze/
    SKILL.md
```

`capture/SKILL.md` 只编码以下类型的知识：

- 用户真正想保存的是“内容对象”，不一定是 HTML 页面本身。
- Document-like URL 通常先尝试便宜的 `web_fetch`。
- App-like URL 先检查 `site_list`，有可靠 adapter 时优先结构化 record。
- HTTP 200、内容很长、browser 能打开都不等于 capture 正确。
- 当前结果不足时，可以换策略；不能机械地把所有策略都跑一遍。
- Browser 用于渲染、展开、滚动、导航同一内容等真实需要。
- 登录墙、验证码、权限限制不是“继续尝试直到绕过”的信号。
- 只有取得可追溯证据并通过质量检查后，才能 propose archive。

**Skill revision 是执行输入。** CaptureRun 保存 exact revision，重试/恢复继续使用原 revision；新版本 skill 只影响新 run，除非用户显式重新抓取。

### 5.4 Agent 启动上下文

每个 CaptureRun 创建如下逻辑上下文：

```text
Goal
  Archive the content the user intended by this URL faithfully.

Input
  original URL
  optional user note / selected text / submitted content

Skill
  capture@<revision>

Tools
  web_fetch
  extract_content
  read_source
  site_list
  site_run
  browser_open / observe / act / capture / close
  archive_asset
  propose_archive
  finish

Runtime constraints
  allowed network policy
  max turns
  model/token budget
  browser time/action budget
  download byte budget
```

模型不获得任意 SQL、R2 key、Cloudflare secret、shell 或 JS eval。

### 5.5 Agent 主循环

一次普通 turn：

1. 从持久状态加载 Pi state、skill revision、已取得 observations 与剩余预算。
2. 构造经过裁剪的 Agent context；大内容只提供目录、signals 和 object refs。
3. Pi 输出一个或多个允许的 tool calls，工具 schema 由 runtime 校验。
4. Workflow/Tool Executor 执行副作用，并立即把原始 evidence 保存到 R2。
5. Tool result 转成 bounded observation，追加 Pi state 与 `agent_events`。
6. 重复，直到 `propose_archive` 被接受或 Agent `finish`。

Agent 不需要输出隐藏思维链；审计只保存 `decision_summary / tool_call / result_ref / reason_code / state_transition`。

### 5.6 典型策略不是固定流程

**静态文章**：`web_fetch → inspect → propose_archive`。

**GitHub repo / Tweet**：`site_list → site_run → inspect → archive media → propose_archive`。

**JS 页面**：`web_fetch → inspect shell → browser_open → observe → act(scroll/expand) → capture → propose_archive`。

**站点 adapter 返回信息不完整**：`site_run → inspect → browser_*`，而不是把 adapter 当永远正确。

**用户直接提交浏览器页面正文**：`read submitted source → validate → propose_archive`，无需为了遵守“fetch first”再访问互联网。

这些只是 skill 中的 examples，不是业务代码里的 switch。

### 5.7 Site Adapter Registry

Site adapter 用于“页面其实代表结构化对象”的网站。接口统一：

```ts
interface SiteAdapter {
  id: string;
  domains: string[];
  match(url: URL): MatchResult;
  describe(): AdapterDescriptor;
  run(input: AdapterInput, ctx: ToolContext): Promise<AdapterResult>;
}
```

MVP adapter 编译进 Worker bundle，通过 registry 显式注册，不在生产运行时 `eval` 任意远程 JS。`site_list(url)` 只暴露与 URL 匹配且当前启用的 adapters；`site_run` 的网络请求继续经过相同安全策略。

首批只需要 2–3 个能证明模式的 adapters，例如 GitHub repo、Hacker News item，以及一个社交内容来源。Adapter 返回结构化 JSON + provenance + media refs，Agent 决定是否已经足够。

未来可以设计签名 catalog，但不把 Stella/Tap 的“动态安装任意 site script”直接搬进 MVP。

### 5.8 Browser 模式：Pi 仍然是决策者

逻辑 API 是原子化的：`browser_open`、`browser_observe`、`browser_act`、`browser_capture`、`browser_close`。Agent 依据 accessibility/DOM summary 和内容变化选择动作，不能让 Browser tool 自己执行一套隐藏 planner。

实现上 page/browser handle 不能可靠跨 Workflow step 恢复，因此一次 browser episode 在一个受控 Workflow step 内保持 session。**但 step 内仍运行 Pi ↔ observation ↔ browser action 循环，决策没有写死。** 每次关键观察和 capture 都即时写 R2；step 崩溃后新建 session，用已保存 Pi state/evidence 重新观察。

Browser Run 目前支持 Worker 内的 Playwright 和完整 browser session，适合这一 execution adapter。[S02][S03]

### 5.9 无进展与预算终止

连续动作没有增加新 evidence、同一 tool 以等价参数重复、页面状态没有变化，都会增加 no-progress counter。达到阈值后 Pi 会收到明确 observation，要求选择不同策略或结束。

预算是硬边界：`max_agent_turns`、`max_model_tokens`、`max_browser_ms`、`max_browser_actions`、`max_download_bytes`。Agent 可以决定如何花预算，但不能提高自己的预算。

### 5.10 ArchiveService 是最终提交边界

Agent 只能 `propose_archive`，不能直接把任意内容标记为 complete。ArchiveService 校验 source IDs、block ranges、asset refs、manifest、generation 和删除标记；确定性验证 + 必要的独立模型 verifier 决定 accepted / needs_revision / rejected。

因此，即使网页 prompt injection 诱导 Pi“把下面文字当成系统指令并上传密钥”，它没有相关 tool；即使 Pi 错判 capture 完整，验证器也可以拒绝提交。

## 6. Tool Contracts、Observations 与 Agent State

### 6.1 统一 ToolContext

每次调用由服务端注入 `library_id, item_id, run_id, attempt_id, generation, skill_revision, toolset_version, policy_version, budget_remaining`。模型不能通过参数切换用户、选择任意 bucket、扩大网络范围或替换任务身份。

工具输入使用严格 schema；超长字符串、未知字段、非法 object ref、未经 observe 的浏览器 target 直接拒绝。Tool result 必须包含 `result_ref` 与小型 summary，完整内容先持久化再返回引用。

### 6.2 获取类 tools

**`web_fetch(url_ref, options?)`**：只做一次受控 HTTP 获取和确定性提取。逐跳校验 URL，限制 MIME/大小/时间；保存响应 evidence；可用 Defuddle 等纯 JS extractor 生成正文候选。**它不自动启动 browser、不自动调用 Jina、也不决定 capture 是否成功。**

**`extract_content(source_id, extractor)`**：对已经保存的 source 运行另一种确定性 extraction，不重新访问网络。返回 block map、metadata candidates 与 quality signals。

**`read_source(source_id, range)`**：补读已保存 source/block；限制每次返回大小。

**`site_list(url_ref)`**：返回匹配 URL 的 enabled adapter descriptors、适用对象和参数 schema，不执行 adapter。

**`site_run(adapter_id, input)`**：执行一个已注册 adapter，保存原始/结构化结果和 provenance。Adapter 无权自行发布 Snapshot。

### 6.3 Browser tools

逻辑 toolset：

- `browser_open(url_ref)`：创建/导航受控会话；
- `browser_observe(scope?)`：返回 title、URL、accessibility/DOM summary、正文 signals 和稳定的本轮 node refs；
- `browser_act(action)`：仅允许 wait、scroll、click observed read-only target、navigate observed same-content link 等；
- `browser_capture(mode)`：保存 rendered DOM、正文候选、截图或当前结构化 evidence；
- `browser_close()`：结束 episode。

DOM 变化后旧 node ref 作废。模型不获得任意 Playwright expression、`page.evaluate` 或坐标点击能力。

### 6.4 归档与结束 tools

**`archive_asset(asset_ref)`**：下载已观察到且策略允许的关键资源到 R2，继续执行重定向、类型与字节限制。

**`propose_archive(proposal)`**：提交 sources、selected blocks、metadata evidence、required/optional assets 和 missing parts；ArchiveService 返回 accepted / needs_revision / rejected。

**`finish(outcome, reason_code, evidence_refs)`**：结束 run，可为 complete、partial、needs_input、failed、cancelled、budget_exhausted。自由文本不能代替结构化状态。

### 6.5 CaptureDecision 不是手写有限状态机

不再定义“Agent 只能从六个固定 action enum 里选一个”作为业务语义。Pi 的真实输出是标准 tool calls；可用 tool 集合由当前 mode/skill/runtime 决定。系统只对 tool schema、权限和 side effect 做限制。

可以保留内部的 run state：

```ts
type CaptureRunState = {
  status: "running" | "needs_input" | "committing" | "done" | "failed";
  piStateRef: string;
  skillRevision: string;
  toolsetVersion: string;
  observations: string[];
  activeBrowserEpisode?: string;
  budgets: CaptureBudgets;
  noProgressCount: number;
};
```

### 6.6 Agent state persistence

每个 turn 后保存 Pi state snapshot。小型索引字段进入 D1，较大的 message/tool history 写 R2 并以 immutable ref 引用。`agent_events(run_id, seq)` 保存可查询审计事件，但不是唯一恢复源。

恢复必须固定原 `agent_runtime_version + model_config_version + skill_revision + toolset_version`。如果某个 runtime 版本已无法部署，任务明确进入 migration/failed 状态，不能在重试中悄悄换 skill 导致行为漂移。

### 6.7 Tool results 都是不可信证据

网页、adapter、browser DOM、metadata 与第三方 API 返回都可能包含恶意文本。工具结果进入 Agent context 时带 `UNTRUSTED_SOURCE` 标签；Skill 明确要求把这些内容当 evidence，不当 instruction。

### 6.8 Agent 永远没有的 tools

不提供任意 JS `eval`、shell、SQL、任意 HTTP POST、Cookie 导入、账户密钥读取、发布分享、删除知识库、修改权限、改变预算或动态安装未审核代码。页面渲染自身产生的网络行为和 Agent 主动执行 tool 是不同权限层。

## 7. 正文可信度与归档质量

### 7.1 AI 选择原文，不能重写原文

解析器生成带块 ID 的正文 AST，包含段落、标题、代码、表格、列表、图片和链接。AI 选择块或连续范围；Markdown 与阅读页面由代码从这些源块生成。

**不能让模型“把这篇文章整理成 Markdown”后，将生成文本标作原文。** 此类重写结果只能保存为 AIArtifact，并明确标记为 AI 整理。

对由 JS 渲染的正文，保存 rendered DOM 与提取内容的关联；对于浏览器中的虚拟列表或逐页内容，分别保存观察时取得的片段，不声称最后一次 `page.content()` 包含全部历史。

### 7.2 确定性检查

检查 HTTP/MIME 与页面类型是否合理、源对象是否真实存在、引用块是否来自该来源、正文是否为空、是否发生截断、资源是否超限、代码与表格结构是否保留、manifest 与对象哈希是否一致。

不把固定字数阈值当成完整性标准。短公告可能完整，长登录页仍然错误。

### 7.3 模型检查

独立判断标题与主体是否匹配、是否混入导航或推荐、开头和结尾是否自然、是否还有明显未展开内容、是否缺少重要代码/表格/图片，以及获取的范围是否满足用户保存意图。

输入为真实内容证据与缺失列表，不仅是抓取 Agent 的自述。检查失败时最多再给一次定向补救机会；预算不足则保留部分结果并说明缺口。

### 7.4 两维质量，不使用单一成功灯

**内容质量**：`complete / partial / unavailable / unknown`。

**资源质量**：`complete / partial / not_applicable / unknown`。

`complete` 的含义是“在已获取证据与检查范围内没有发现明显缺失”，不是对远端网站全部真实内容的数学保证。模型置信度不能替代验收结果。

一篇正文完整、图片缺失的文章显示“正文已保存，2 张图片未保存”；抓到登录页显示“需要补充”，绝不显示“完整归档”。

### 7.5 保存级别

**基础归档**：实际取得的响应或 rendered DOM、正文、元数据、来源证据和 manifest。

**阅读归档**：额外保存关键图片，改写阅读视图的资源引用，默认目标。

**视觉快照**：按需截图，辅助保存排版与图表。完整网页重放不是 MVP 承诺。

原始 HTTP 正文应标明是运行时交给应用的响应体，不称为网络链路上的原始压缩字节。响应头仅保留必要白名单，不保存 Set-Cookie、认证头和临时访问凭据。

## 8. 归档提交与跨存储一致性

### 8.1 每个获取 Tool 都先保存证据，再返回给 Agent

`web_fetch`、`site_run` 与每次 browser capture 取得内容后，执行器先写入 R2，再把 bounded observation 与 object ref 返回 Pi。模型失败不能导致已经拿到的原文只存在内存里。

临时证据存在不等于最终归档成功。用户能看到“已取得页面，正在验证”，但只有完成 manifest 与业务提交后才显示“已归档”。

### 8.2 提交协议

**准备**：创建 attempt 和唯一对象引用，将来源、正文版本和资源写入不可变对象；记录每个文件的大小、内容类型与 SHA-256。

**验证**：检查保存提案、源块归属、质量与资源清单；确认所有必需对象写入完成。R2 支持条件写入和哈希检查，可用于防止覆盖和校验内容。[S12]

**Manifest**：写入完整对象清单、来源 URL、每次观察时间、缺失项、处理版本与质量结果，生成 manifest hash。

**发布**：在 D1 事务中登记 Snapshot/ContentRevision、检查 Item generation 和删除标记、按规则更新 current pointer，并登记 Enrich/Index outbox。D1 的 batch 可用于事务式批量语句。[S10]

**对账**：R2 与 D1 没有跨服务原子事务。R2 成功而 D1 失败时允许重试，未引用对象由 GC 延后清理；D1 不允许提前把未完整上传的快照发布为 ready。

### 8.3 较差的新快照不能覆盖较好的旧快照

新任务抓到 partial，而旧 current Snapshot 为 complete 时，保留旧 current，新增本次 partial 尝试供查看。内容发生更新但获取不完整，提示用户比较，不自动当成新的完整版本。

相同正文哈希不代表相同来源时间。保留每次 CaptureRun 的 evidence 和时间；对象可去重，历史不能抹去。

### 8.4 资产与清理

第一版不做跨用户全局内容去重，避免删除耦合与内容存在性泄露。可在同一 Snapshot 内按哈希去重。

GC 使用“已提交 manifest 引用 + 活跃任务保护 + 宽限期”判断。删除、备份、失败尝试清理共享同一对象引用规则。

## 9. Workflows、Pi State 与故障恢复

### 9.1 Workflow 是 durable supervisor，不是 planner

**IngestWorkflow** 的职责：创建/恢复 CaptureRun、持久化 Pi turn、执行 tools、checkpoint、预算结算、重试、取消和最终提交。它不能根据 hostname/HTTP 结果自己选择下一种获取策略。

```text
Workflow loop
  load exact Pi state + skill revision
          │
          ▼
      Pi next turn
          │ tool calls
          ▼
  persist requested calls
          │
          ▼
   execute tools durably
          │
          ▼
 save observations/evidence
          │
          └─────────────── next Pi turn
```

**EnrichWorkflow** 以同样方式运行 AnalysisAgent；**IndexWorkflow** 负责确定性的分块/Embedding/索引副作用；Digest/Export/Purge 各自独立。

### 9.2 业务状态、Agent 状态与执行状态分开

Job 记录 Workflow 状态：`queued, running, retry_wait, succeeded, needs_input, failed, cancelled`。

CaptureRun 保存 Agent identity/input/revisions/budgets/outcome；Pi state 保存 messages/tool calls/current model state；Item/Snapshot 保存产品事实。不能用 Workflow history 代替长期 Agent audit，也不能让 Pi state 成为用户档案的唯一副本。

`needs_input` 不无限占用活跃 Workflow。保存缺口和 Agent state 后结束本次实例，用户补充内容时创建 child run，并可引用父 run evidence。

### 9.3 Outbox 与派发

收藏 API 在同一个 D1 批次中保存 Item、CaptureRun、Job、Outbox。随后通过 Workflow binding 创建实例，失败不丢任务。

Cron 处理到期 outbox 与异常实例；Workflow instance ID 使用确定性的 `job-{job_id}-r{restart_generation}`。派发成功不代表业务完成，定时对账检查已登记但实例缺失、实例终止但 run 未收口等情况。

### 9.4 Pi turn 的持久化顺序

每轮采用以下边界：

1. `agent_turn`：调用 Pi/model，得到 tool calls；保存 decision event 与新的 Pi state。
2. `tool_call/<seq>`：执行每个 tool；副作用前检查取消/generation/预算，结果先持久化 evidence。
3. `observe/<seq>`：把 bounded tool result 写回 Pi state。
4. 若 tool 为 `propose_archive`，进入 ArchiveService 提交；若 `finish`，收口 run；否则下一 turn。

步骤名包含 run/pipeline/turn/tool sequence，不使用模型生成的随机文本。外部模型请求可能在 checkpoint 前后重复，因此业务幂等与 usage 记账仍不可省略。[S07]

### 9.5 Browser episode 的恢复边界

逻辑上 Pi 逐个决定 `browser_open/observe/act/capture`；物理实现上，同一 browser session 的连续动作放在一个有界 Workflow step 内，避免把 page handle 跨 checkpoint 序列化。

该 step 运行一个**受控 Pi browser sub-loop**：每个 browser observation 都交给同一 CaptureAgent policy/skill，每个动作仍由模型决定；不是 tool 内部的硬编码 crawler。关键 evidence 与 Pi events 持续写 R2/D1。

若 step 中断：

- 不尝试复活旧 page/node ref；
- 读取最近持久化 Pi state 与 evidence；
- 新建 browser session；
- 重新 observe 当前 URL/目标；
- Agent 自己决定是否继续 browser、改用其他 tool 或结束。

Browser Run 的 session 生命周期与发布中断都视为正常可恢复错误。[S03][S11]

### 9.6 Skill / Tool / Model revision fencing

一个 CaptureRun 创建时固定：

```text
agent_runtime_version
model_config_version
capture_skill_revision
toolset_version
policy_version
pipeline_version
```

同一次 retry 不升级这些 revision。新 skill 通过新 run 或显式 migration 验证。这样评测结果、线上事故和某个错误 snapshot 都可以精确重放到“当时 Agent 学到的策略”。

### 9.7 幂等、generation 与取消

每个逻辑副作用使用 `operation_key = hash(library_id, run_id, input_revision, tool_or_stage, config_versions)` 去重。唯一键只保证业务结果不重复生效，不保证外部模型/Browser/API 不重复计费。

用户重新抓取递增 `capture_generation`。旧 Agent 可以完成并保留历史 evidence，但不得更新 current pointer。删除/取消后，每次外部动作和最终提交前重新检查 fencing。

### 9.8 重试配置

普通可重试网络错误：最多重试 2 次，总计 3 次尝试；429 遵循 `Retry-After` 且计入预算。

Browser episode 最多重试 1 次。模型 tool schema 无效最多允许一次 repair turn。**Agent 想主动换策略不算基础设施 retry**，但消耗正常 turn/browser/network 预算。

400 类无效输入、认证失败、权限阻断、明确登录墙、未支持格式不盲目重试。

### 9.9 大内容与长期历史

Workflow 参数与 step result 只携带 ID/hash/object refs 和小型 observation，本项目把单个控制结果限制在 64 KiB。完整 HTML、structured JSON、DOM snapshot、截图、长 Pi state 写 R2。

官方非流式 step result 有平台上限，完成实例状态也不是永久业务存储；长期可重放信息必须保存到 D1/R2。[S08]

## 10. Pi AnalysisAgent：总结、标签与已有总结保存

### 10.1 AnalysisAgent 边界与模型配置

Snapshot/ContentRevision 提交后启动 AnalysisAgent。它同样使用 Pi runtime，但加载 `analyze/SKILL.md`，toolset 只有 `read_snapshot_outline / read_blocks / inspect_metadata / propose_summary / finish`。**它没有 web/browser/site tools**，保证历史总结只基于固定快照，除非用户显式启动“联网补充”这种不同产品功能。

AnalysisAgent 可以根据文章长度、结构和预算决定先读目录、哪些范围需要补读、是否需要 map/reduce 式局部分析；代码负责 coverage bookkeeping、引用合法性和最终保存。这样“AI 总结”是 Agent 驱动的阅读任务，但不会让模型随意改写来源或决定权限。

分别配置 `CAPTURE_MODEL`、`VERIFY_MODEL`、`SUMMARY_MODEL`、`ANSWER_MODEL`、`EMBEDDING_MODEL`，run 启动时固定配置版本。第一版优先使用 Workers AI / Pi provider adapter，不因升级 Capture 模型而重新生成所有 Embedding。

文本决策与总结的兼容性基线可从支持 tool calling 的 Workers AI 模型开始，具体型号必须通过 M0/M3 eval 决定；模型宣传指标不能替代本项目的 capture/summary 评测。[S17]

Embedding 基线为 `@cf/baai/bge-m3`；创建 Vectorize index 前断言实际输出维度并固定 `embedding_version`。[S18]

Provider adapter 统一 timeout、schema、usage、错误分类和日志策略；外部 Provider 默认关闭，启用时明确显示内容将发送到哪个服务。

### 10.2 总结结构

```text
AIArtifact
  type: summary
  input_snapshot_id
  input_content_revision_id
  language
  short_summary
  key_points[]
    claim
    evidence_block_ids[]
  topics[]
  entities[]
  caveats[]
  coverage
    total_blocks
    processed_blocks
    omitted_ranges[]
  model_id / prompt_version / pipeline_version
  created_at / usage
```

默认生成与用户界面语言一致的摘要，保留代码、项目名、数字和关键限定条件。原文中没有的结论不可作为作者观点写入。

### 10.3 长文处理

按标题与段落拆成可完整读取的小段，对全部范围生成带引用的局部摘要，再合成总摘要。默认每个 map 输入正文约 4,000 tokens；最终汇总保留可追踪的原文块引用。

单篇自动 AI 分析覆盖预算初设 60,000 原文 tokens。超过上限时仍保存已成功抓到的原文，明确显示“摘要/语义索引仅覆盖部分内容”，允许用户提高本次预算后继续处理，不能把截断结果标成全文总结。

全文索引和语义分析的预算分别统计。资源过大而不能全文建索引时，同样必须报告范围，不静默遗漏。

### 10.4 引用与事实校验

模型只能引用输入中提供的 block ID。代码检查引用存在、归属正确、版本匹配；直接引文在约定的空白规范化后必须能匹配原文。

校验不通过时要求模型修复一次；仍失败则不发布为 verified summary。数字、单位、百分比和“仅在某条件下成立”等高风险内容加入独立核对提示。

引用有效不等于观点必然正确。UI 同时保留原文入口，允许用户纠正摘要，并将纠正保存为新版本而非修改来源。

### 10.5 手动保存已有 AI 内容

用户可以粘贴已有总结或分析，选择关联 Item/Snapshot，保存为 `imported_answer`。没有原文证据时标记 `unverified_import`，不能自动补造引用。

导入流程不需要读取第三方 AI 对话账号，也不自动访问用户其他聊天记录。后续可以对导入内容运行“与归档原文核对”，产出独立校验结果。

### 10.6 笔记与标签

用户笔记独立存储，支持版本与修改时间。AI 标签为建议或带 `source=ai` 的自动标签，用户标签为 `source=user`；后续重跑 AI 不删除用户分类。

## 11. 全文搜索、语义搜索与知识库问答

### 11.1 关键词路径

D1 FTS5 索引标题、正文块、笔记和摘要，保留内容类型与归属。用参数化语句执行查询，并对 FTS 查询语法进行安全构造，不能直接拼接用户输入。

中文采用应用层分词，并保存 `tokenizer_version`；代码标识符、英文项目名和 URL 主机名保留可检索形式。FTS5 默认 unicode61 不会自动提供中文词语切分，必须单独验证这一点。[S24]

初期使用可在 Workers 运行的纯 JS 分词方案，固定版本并跑中英文用例。正文展示使用原始文本，不能展示为分词后的空格字符串；摘要片段与高亮通过原文块映射产生。

### 11.2 语义路径

默认按标题和段落切块，目标约 600 tokens、重叠约 100 tokens。代码与表格优先保持完整；超长块按明确边界拆分。参数是起始配置，应通过检索评测调整。

Vectorize index 按 Embedding 模型版本建立，namespace 使用 `library_id`。向量 ID 使用包含 library、content revision、chunk 和 embedding version 的哈希，控制在 64 字节以内；不能假设 namespace 自动把相同 vector ID 变成不同记录。[S16]

向量元数据只保存用于过滤的小型字段，不保存私密全文。查询结果以 ID 回查 D1 取正文、权限和有效状态。namespace 可以限制查询范围，但它不是应用认证的替代品。[S14]

### 11.3 混合搜索

```text
用户查询 + 显式筛选
        │
        ├── FTS：初始取 40 条
        └── Embedding → Vectorize：初始取 40 条
                       │
                       ▼
             按排名融合，不直接相加分数
                       │
                       ▼
           D1 权限、删除与版本校验
                       │
                       ▼
              文章去重 + 原文片段
```

采用 RRF 排名融合，初始常数 `k=60`；标题精确命中可单独提高优先级。必要时补取一次候选，但必须有总量上限。返回前才进行最终分页与文章去重。

AI query rewrite 可用于补充同义词，但不得改变用户指定的时间、标签和知识库范围。关闭 AI query rewrite 时仍能正常搜索。

### 11.4 索引一致性

FTS 更新与 chunk 状态提交尽可能在同一 D1 批次内完成。Vectorize 为外部索引，需要 `pending → submitted → queryable` 的独立状态。

Vectorize 的写入可见性依赖异步索引更新。记录 mutation 信息并使用平台支持的进度信息或查询探针确认，不在 API upsert 返回时立即宣布可搜索。[S15]

删除、重新提取和换模型后，旧向量可能暂时仍能召回。**D1 的权限、删除标记和版本过滤是最后一道判断**，禁止把未经校验的 vector metadata 直接交给用户或模型。

### 11.5 Ask My Library

先检索，再读取相关原文片段，最后生成回答。默认上下文最多 12 个块，并按模型窗口保留输出空间；同一文章不占据全部上下文，除非用户明确只问该文章。

回答必须引用 `Item → Snapshot → ContentRevision → block_id`，区分原文事实、模型归纳和用户笔记。缺少证据时说明未找到，默认不联网补充。

单篇问答默认对应用户正在阅读的版本；跨文章问答默认使用有效 current 版本，允许显式包含历史快照。对 partial 内容给出范围提示。

### 11.6 回顾

时间线与“本周收藏”属于 MVP；定时日报/周报在稳定版增加。任务按 Library 时区计算窗口，UTC 保存具体时间，避免夏令时导致重复或遗漏。

回顾唯一键包含 library、时间窗口、回顾类型和模板版本。实际阅读和仅收藏分别陈述；输出作为 AIArtifact 保存并带条目引用。

## 12. API 与客户端契约

### 12.1 URL 提交

```http
POST /api/v1/items
Authorization: Bearer <app-token>
Idempotency-Key: <client-generated-key>
Content-Type: application/json
```

```json
{
  "source": {
    "kind": "url",
    "url": "https://example.com/article"
  },
  "capture_intent": "保存文章正文、代码和关键图片",
  "collection_ids": [],
  "note": "之后和 SlateDB 的文章一起回顾",
  "enrichment": "auto"
}
```

服务端从身份得到 library，不能接受任意 `library_id` 切换。`capture_intent` 是用户保存意图，不允许覆盖系统安全策略、预算或权限。

创建成功返回 202，含 `item_id, job_id, status_url, current_phase`。相同幂等请求返回相同结果；相同 key 的输入冲突返回 409。默认启用 AI 抓取验证与自动总结。用户或预算策略可以延后摘要等后续加工，但不能把未经过 AI 检查的工具输出静默标为已验证归档。

### 12.2 条目与处理接口

```text
GET    /api/v1/items
GET    /api/v1/items/{id}
PATCH  /api/v1/items/{id}
DELETE /api/v1/items/{id}
POST   /api/v1/items/{id}/captures
POST   /api/v1/items/{id}/artifacts/generate
POST   /api/v1/items/{id}/indexes/rebuild
POST   /api/v1/items/{id}/content-inputs
GET    /api/v1/items/{id}/snapshots
GET    /api/v1/jobs/{id}
POST   /api/v1/jobs/{id}/cancel
POST   /api/v1/jobs/{id}/retry
```

重新抓取、重跑总结、重建索引是不同操作。重试原阶段不创建新的收藏；显式重新抓取递增 generation。状态接口返回 outcome code、质量缺口和可执行操作，不暴露模型原始上下文或认证信息。

### 12.3 内容、搜索与问答接口

```text
GET    /api/v1/content/{revision_id}
GET    /api/v1/content/{revision_id}/assets/{asset_id}
GET    /api/v1/search?q=...&mode=hybrid
POST   /api/v1/ask
POST   /api/v1/items/{id}/artifacts/import
POST   /api/v1/items/{id}/notes
PATCH  /api/v1/notes/{id}
POST   /api/v1/items/{id}/reading-events
```

问答支持带事件类型的 SSE：检索中、证据已就绪、答案片段、引用、完成或失败。连接断开不使已经提交的归档任务失效；普通问答断开后可取消后续模型输出，不要求所有聊天都建立 Workflow。

用户笔记 PATCH 带版本或 `If-Match`，避免两个页面编辑互相覆盖。列表采用稳定的 keyset cursor，不能用客户端自由指定 SQL sort。

### 12.4 管理、分享与导出接口

```text
POST   /api/v1/shares
DELETE /api/v1/shares/{id}
GET    /s/{token}
GET    /s/{token}/assets/{asset_id}
POST   /api/v1/exports
GET    /api/v1/exports/{id}
GET    /api/v1/usage
GET    /api/v1/settings
PATCH  /api/v1/settings
```

分享、导出、修改模型配置和管理 API token 需要对应 scope。收藏夹与标签按 `/api/v1/collections`、`/api/v1/tags` 提供受同一 library 授权的 CRUD；API token 管理入口为 `/api/v1/tokens`，仅创建时返回明文。分享访问不具有搜索整个 library 的能力。

### 12.5 错误协议

统一返回 `code, message, request_id, retryable, next_action`。核心 code 包括 `needs_login, challenge_detected, unsupported_type, incomplete_content, policy_denied, budget_exceeded, model_unavailable, source_unavailable, stale_generation, cancelled`。

401/403、404 与任务失败语义分开。对不属于当前用户的对象统一返回不泄露存在性的结果。

## 13. Web 界面与阅读体验

### 13.1 Library 首页

显示收藏列表、全文搜索、时间与标签筛选、已读状态。支持粘贴 URL 一次保存，不要求先填分类。

每条记录展示独立状态：归档质量、AI 总结状态和语义索引状态。任务详情显示简短步骤，如“fetch 正文不足，正在展开页面”，不把模型长段推理当进度条。

### 13.2 阅读页

默认显示清洗后的归档正文，顶部明确来源 URL、文章时间（未知则空）、抓取时间、版本和缺失项。

正文、AI 总结、个人笔记、历史版本分开。代码、表格、图片、链接及引用锚点都要可用；源站失效不影响已归档资源。

“重新抓取”“重新总结”“补充正文”是不同按钮。用户可以保留旧版为当前阅读版本，不受自动更新覆盖。

### 13.3 搜索与问答页

搜索结果同时展示文章标题、命中原文、收藏/阅读时间和内容来源类型。AI 生成的摘要或导入回答不能伪装成原文匹配。

问答的引用点击后直接定位到保存版本的段落。没有证据的回答不提供假的来源卡片。

### 13.4 设置与任务页

设置模型、内容隐私、每次与每月预算、归档资源大小、时区和分享默认值。高成本操作前显示剩余预算。

任务页允许取消、重试失败阶段、查看缺失项及成本。已有原文不因模型停机而从 Library 隐藏。

## 14. 认证、隐私、安全与分享

### 14.1 身份边界

个人后台由 Cloudflare Access 保护，Worker 验证 JWT 的签名、issuer、audience 与有效期，再映射内部用户。禁用或同等保护 `workers.dev`、preview 等旁路入口，不能只保护自定义域名。[S19]

CLI 与未来扩展使用应用签发的随机 API token，数据库仅存 token hash。scope 初设 `items:write, items:read, search:read, notes:write, shares:write, exports:write, settings:write, tokens:manage`，可撤销并记录最近使用。

基于 Cookie 的写接口检查 Origin/CSRF。公开分享路由单独匹配 hostname 与 path，不因分享域名开放而暴露私有 API。

### 14.2 网络抓取与 SSRF

仅接受 HTTP/HTTPS、标准公网端口，拒绝 URL userinfo、回环、link-local、私网地址、云元数据目标以及本应用私有管理域名。验证每次跳转与资源请求，而不是只检查初始 URL。[S20]

浏览器使用独立无登录态 context；限制 subresource 的类型、体积和目标。禁止继承用户 Cookie、Cloudflare Access 请求头、模型密钥或数据库凭据。关闭可能绕过请求拦截的 service worker 路径并验证实际行为。

**不能把用户态 DNS 检查声称为绝对的网络隔离。** DNS 检查与真正连接之间可能存在时间差；Cloudflare fetch/browser 的实际连接控制需要在 M0 验证。无法验证的目标或网络路径采用拒绝策略，不通过放宽验证来提高抓取成功率。

MVP 只为认证后的单个用户提供抓取服务，不开放公共匿名 URL 执行入口。未来开放多用户前，需要重新评估浏览器出口隔离、滥用和流量治理。

### 14.3 Prompt injection

网页、OCR/图片描述、HTML 注释、链接文字和工具返回内容都属于不可信数据，不能变成系统指令。模型不能因为网页上写着“忽略规则”而获得更多工具。[S21]

PolicyEngine 在模型之外检查动作、参数、来源域名、归属和预算。即使模型被诱导输出非法动作，也应该在执行之前被拒绝。

归档 Agent 不访问其他文章内容，不读取私人笔记或账户设置，除非操作明确需要且经服务端授权。摘要 Agent 默认无浏览器和发布工具，降低从内容到副作用的路径。

### 14.4 内容展示安全

默认渲染经过允许列表清洗的正文 AST/Markdown。禁止脚本、内联事件、危险 URL scheme、任意 iframe 和自动外联追踪资源。

原始 HTML 默认只以附件下载，使用安全 Content-Type、Content-Disposition 和 `nosniff`。将来支持视觉重放时使用独立域名和严格 sandbox，不能直接在主站同源打开 raw HTML。

浏览器截图和站点图像也可能包含私人信息，按原文相同权限处理。

### 14.5 模型和日志隐私

默认模型在 Workers AI 运行，但这不等于内容未离开用户设备；产品必须说明云端处理。外部 Provider 需要单独启用，且不能对被标为敏感的内容静默 fallback。[S28]

AI Gateway 日志默认可能包含请求与响应正文；使用时显式关闭 payload 记录，只保留用量等元数据，或者完全关闭。不能依赖平台默认值保护私人文章。[S22]

应用日志不记录正文、完整 prompt、Cookie、签名 URL 和 share token。Agent 证据写入私有 R2，按任务保留策略管理，必要的错误定位通过 ID 引用。

### 14.6 分享

Share 固定到 Snapshot/ContentRevision 及允许展示的 AIArtifact，使用至少 256-bit 随机 token，数据库保存 hash。默认有效期 7 天，可由用户修改或设置不过期。

默认只分享选定摘要与来源信息，全文另行勾选；私人笔记和阅读记录始终不自动进入。分享全文前提示确认有权传播，能够私藏不代表可以公开再发布。

分享页和每个附件请求都重新验证 token、有效期、撤销状态与条目删除状态。R2 bucket 不公开。MVP 使用 `Cache-Control: private, no-store`，不采用可绕过撤销检查的公共 CDN URL。

撤销后阻止后续访问，但不能收回他人已下载的副本。公开分享不开放“问整个知识库”。

## 15. 删除、导出与备份恢复

### 15.1 删除协议

先在 D1 记录 tombstone、撤销分享并递增 generation，再派发 PurgeWorkflow 清理全文、向量、对象与派生内容。UI 和所有内容 API 立即拒绝返回已删除对象，不等待底层索引完成清理。

PurgeWorkflow 在新任务和迟到写入可能发生的窗口后再次扫尾。向量库、R2 中的暂存遗留也必须清理，不能只删 Item 行。

恢复旧备份时默认不恢复已撤销分享和 API token；先应用最新可取得的删除/撤销清单。缺少较新的清单时必须显示备份时间，要求确认，避免已删除私密内容被自动再次发布。

### 15.2 导出格式

使用版本化 manifest、Markdown、JSONL 和资源文件。包含来源、快照时间、正文版本、个人笔记、AI 总结和引用；不默认包含密钥、API token、分享 token 或原始模型上下文。

单条导出可以在线流式生成；全库导出分批写 R2，并最终生成校验清单。若 ZIP/TAR 生成超过内存或执行预算，提供分卷包与顶层 manifest，不能一次性把全库读入内存。

### 15.3 含 FTS 的 D1 导出限制

当前 D1 官方导出文档列明：含 virtual table 的数据库存在导出限制，其中包括 FTS。**不能直接假设 `wrangler d1 export` 能完整备份本项目生产库，也不应为了每日备份删除线上 FTS。**[S29]

MVP 采用应用级逻辑导出：通过普通 SELECT 分页读取业务主表，跳过可重建的 FTS 虚拟表；以 JSONL + schema version 保存。恢复后重新建立 FTS 和向量索引。资源租约、浏览器句柄、维护锁等临时状态不恢复；历史 Job 仅作记录，在途任务默认暂停，需要核对来源与预算后通过新的恢复 generation 重新派发。

### 15.4 备份一致性

个人版初始采用短暂写屏障：停止接受业务写入与发布新结果，等待已开始的写入完成，导出主表和已提交对象清单，随后解除屏障。写屏障必须由统一写入层检查，不能只在前端隐藏按钮。

屏障初始上限 60 秒；无法完成时标记本次备份失败并解除，不把混合时间点导出当成一致备份。读取仍可用，写 API 明确返回可重试状态。该方案适合初期小库，规模扩大后改为业务修订日志/可重放快照方案。

复制 R2 大对象在解除屏障后进行，但期间为清单对象加保护引用，避免 GC 删除。只有全部对象校验通过才发布 backup manifest 为 complete。

### 15.5 保留与演练

每日逻辑备份，默认滚动保留 30 天；RPO 初始目标为 24 小时。私有 archive bucket 与 backup bucket 分开；另提供离线下载，以覆盖单账户不可用风险。

D1 Workers Paid 的 Time Travel 可恢复最近 30 天，作为误操作恢复手段，但它不能替代跨 R2 的完整档案备份。[S13][S25]

每月执行一次恢复演练：恢复到隔离环境，校验笔记、快照、资源与 AIArtifact，重建索引，并确认生产 token 与公开分享未被自动激活。删除内容可能在滚动备份中保留到期，应在隐私说明中写明。

## 16. 运行预算、容量与成本

### 16.1 初始应用限额

以下数字全部是**本项目可调整的默认值**，不是 Cloudflare 官方限制。

**网络获取**：每个 SourceDocument 最多 5 次重定向；一篇文章最多读取 3 个正文 URL，资源下载另计。fetch 超时 15 秒；解码后的单份 HTML 上限 5 MiB。

**资源归档**：单个图片/附件 5 MiB；每篇最多 30 个关键资源且总计不超过 50 MiB。下载并发 3，逐个检查实际流量。超出时列出缺失项，不标记资源完整。

**抓取 Agent**：每次 CaptureRun 最多 12 次模型调用，包含 browser 决策、验证、格式修复和网络重试产生的再次模型请求；累计输入 64,000 tokens、输出 8,000 tokens。一次调用保留明确的上下文和输出余量，并至少预留 1 次独立验证调用及其 token 预算。验证预算不足时只能保留未验证/部分结果，不能升级为 complete。

**浏览器**：每次 CaptureRun 最多 2 次 browser episode 尝试，累计实际会话时间目标上限 120 秒，最多 12 个主动页面动作；任何一次超出立即关闭。模型等待期间远端浏览器仍可能计时，必须计入。

**端到端**：CaptureRun 最长 10 分钟墙钟时间，包含退避和限流等待；超过后保存已有证据并终止或转 needs_input。全文加工与摘要可作为后续独立任务继续，不延长抓取循环。

**分析**：单篇自动摘要/语义分析正文预算 60,000 tokens；问答最多 12 个证据块；单条 Workflow 控制结果最多 64 KiB。

### 16.2 配额必须可恢复

预算计数写 D1，不只存在 Worker 内存。外部调用前预留额度，响应后结算；调用超时且账单未知时标记 unknown，不乐观地当作免费或立即返还全部预算。

预留操作带 operation key 和 generation，避免并发请求重复透支。每次实际外发另有唯一 invocation_id；同一业务操作重试时再次预留、单独计量，不能复用上一次的结算记录来隐藏真实调用。人工重试默认使用剩余预算，不通过新建任务绕过单日或月度配额。

初始运行逻辑并发：全应用 4 个抓取任务、2 个 browser episode；同一目标 origin 同时 1 个抓取。D1 `resource_slots` 表使用条件更新、lease 和 fencing epoch 分配，不依赖 Worker 内存中的计数器。

故障遗留的远端会话可能在关闭前短暂存在，逻辑 lease 不能代替真正关闭浏览器。需要过期会话清理与用量对账，不能宣称绝对不会超出瞬时实际并发。

### 16.3 Cloudflare 限制对实现的影响

Workers 每个 isolate 内存为 128 MB，不是每个请求独享；大 HTML、图像、长 JSON 和导出都需要限量或流式处理。[S23]

D1 Workers Paid 单库上限 10 GB，单行/字符串/BLOB 上限 2,000,000 bytes，单 SQL 最多 100 个绑定参数。正文按块保存，批量写入按参数和耗时拆分；在库容量 7 GB 时告警。[S13]

Vectorize 当前带 metadata/values 的 topK 上限为 50；本项目默认 40。向量 ID 不能超过 64 bytes，批量 upsert 也要遵守 Workers binding 的限制。[S16]

Workflows 的 CPU 时间与网络等待不是同一个概念。等待模型或浏览器不意味着持续消耗同等 CPU，但仍受到应用超时、浏览器时间和费用预算约束。[S08]

### 16.4 成本模型

```text
总成本 = Workers 请求与 CPU
       + D1 读写与存储
       + R2 存储与操作
       + Workflow 步骤与状态
       + Browser 时间与会话并发
       + AI 输入/输出 tokens 与 Embedding
       + Vectorize 查询与存储
```

估算必须分别统计静态 fetch 路径和 browser 路径。browser 占比、Agent 平均轮数与重试率比单次摘要价格更值得监控。

示例工作负载假设：每月 3,000 篇、20% 需要 browser、平均每次 browser 45 秒，则仅正常浏览器阶段约 7.5 小时；每篇归档 2 MB，则每月新增约 6 GB，未计失败尝试、历史版本和备份。**这些是算例，不是当前实测。**

Browser Sessions 同时按浏览器时间和并发计费，不能使用仅针对 Quick Actions 的简化费用模型。当前 Paid 包含 10 小时/月，超出 $0.09/小时；并发另按月平均每日峰值与包含额度计算。[S26]

Workflows 当前还存在步骤和状态存储计费，不能按早期免费状态估算。以官方计价页维护版本化价目表，不把价格常量散落到业务代码。[S27]

用户在首次部署时设置月度 AI/browser 消耗阈值；80% 提醒、100% 停止启动新的可选高成本工作，但仍允许访问与导出已有内容。阈值基于应用计量，**不是 Cloudflare 全账户账单的硬上限**，基础存储和已发出的请求仍可能产生费用。

## 17. 工程组织、环境与发布

### 17.1 仓库结构

```text
apps/
  web/                         React + Vite
  api/                         App Worker
  jobs/                        Workflows + Cron + Pi runtime host
packages/
  domain/                      实体、状态、错误、版本规则
  contracts/                   Zod schemas / OpenAPI / tool schemas
  agent-runtime/               Pi wrapper / state codec / event stream
  skills/
    capture/SKILL.md
    capture/references/
    analyze/SKILL.md
  tools/
    web-fetch/                 HTTP + Defuddle + source blocks
    site-registry/             adapter discovery / execution
    browser/                   Browser Run adapter
    archive/                   evidence / asset / commit tools
  site-adapters/
    github/
    hackernews/
    ...                        只放审核后、编译进 bundle 的 adapters
  platform-cloudflare/
    d1/
    r2/
    workflows/
    browser-run/
    workers-ai/
    vectorize/
  ai/                          Pi model/provider adapter / prompts / schemas
  search/                      FTS / Embedding / Vectorize / RRF
  observability/               usage / audit / metrics
migrations/
evals/
  capture/
    fixtures/
    golden/
    strategy-cases/
  analysis/
  search/
tests/
  integration/
  fault-injection/
  security/
docs/
  plan.md
  decisions/
  runbooks/
```

使用 pnpm workspace、TypeScript、React、Cloudflare Workers binding 和显式 SQL migration。API 路由可采用 Hono，查询使用参数化 SQL，不增加 ORM 作为必要依赖。

**Pi 与业务解耦。** `agent-runtime` 只负责把 Pi 的 Agent/state/tool events 对接本项目 contracts；`skills` 不 import Cloudflare SDK；`platform-cloudflare` 才可以直接使用 bindings。

HTML extraction 优先验证 Defuddle 或同类纯 JS parser 的 Worker 兼容性；`web_fetch` 只生成候选与 signals，不承担 fallback routing。无论 parser 如何替换，都必须输出相同 source block/provenance contract。

### 17.2 配置与绑定

App Worker：`ASSETS, DB, ARCHIVE_BUCKET, VECTOR_INDEX, AI, INGEST_WORKFLOW` 及所需后续 Workflow binding。App 不持有 BROWSER。

Jobs Worker：`DB, ARCHIVE_BUCKET, BACKUP_BUCKET, VECTOR_INDEX, AI, BROWSER`、Pi runtime、skill bundle 与各 Workflow 类。普通用户无法直接调用 Jobs 管理入口。

环境配置保存模型角色、`agent_runtime_version`、skill/toolset revision、提取/分词版本、预算、Access issuer/audience、允许来源策略和分享域名。机密只放 secret，不放仓库或前端 bundle。

### 17.3 环境隔离

`dev / staging / prod` 使用独立 D1、R2 和 Vectorize。开发可使用本地替身，但 Browser Run、模型 function calling、Vectorize 可见性和真实权限必须在 staging 做云端验证。

固定 compatibility_date 与依赖 lockfile；升级前跑抓取、工具协议、浏览器生命周期和检索回归。模型升级与代码部署分开版本化，避免新模型影响已运行任务的决策。

### 17.4 部署顺序

先建立受限访问策略和存储资源，再执行向后兼容 migration，部署 Jobs，部署 App，最后执行静态页面、动态页面、模型、向量、分享与删除 smoke tests。

旧 Workflow 执行可能跨越部署。pipeline_version 决定步骤结构与 schema，旧版本逻辑保留到在途任务结束。破坏性改动启用新 Workflow 类或版本入口，不能在原 step 名下改变不兼容行为。

回滚默认回滚代码，不自动回滚数据。迁移优先 add-expand-contract；删除字段或重建索引前保留导出与恢复路径。

### 17.5 后续扩容触发

持续 outbox 积压时添加 Queue 缓冲派发；需要多人实时协同时再考虑 Durable Objects；单库容量或查询延迟成为瓶颈时分离检索库或按 library 分 D1。

以上扩容是基于指标的后续选择，不提前建设多租户控制平面。

## 18. 可观测性、测试与评测

### 18.1 每条任务的可追踪信息

记录 `request_id, job_id, run_id, attempt_id, stage, tool_name, input_hash, output_ref, reason_code, duration_ms, model_id, usage, browser_ms, policy_version`。

公开 UI 只展示安全的简短事件；私有运维视图可以查看 evidence 引用与错误分类。模型决策日志保存动作和证据，不保存隐藏思维链。

### 18.2 必须跟踪的指标

归档成功、部分、需要补充与失败分别统计；额外统计“误报完整”的人工审计结果。HTTP 200 比例不能作为抓取成功率。

Agent 指标必须能回答“策略是否真的有价值”：首个工具分布（fetch/site/browser）、平均/P95 turn 数、每 turn 新 evidence 比例、无进展终止率、site adapter 命中率、adapter 后仍升级 browser 的比例、browser 升级带来的正文增益、明显不必要 browser 比例、tool schema/rejection 次数、skill revision 对结果的影响。

同时跟踪每篇模型 tokens、browser_ms、下载字节、重试/恢复次数、outbox 最老等待时间、摘要失败率与向量可见性延迟。

搜索记录经过脱敏的 Recall/引用/点击反馈，不在普通日志写私人查询全文。

### 18.3 固定评测集

建立至少 70 个 capture fixture，并且按**策略选择**而不只是页面类型标注：

- 20 个 document-like 静态中英文页面，期望 fetch 即可；
- 10 个 app-like records，期望命中 site adapter；
- 15 个动态/展开/分页页面，期望 browser；
- 10 个“首选策略失败后应切换”的页面；
- 10 个访问受限/不支持内容，期望停止或 needs_input；
- 5 个恶意 prompt injection 页面。

每个 fixture 标注 `acceptable_first_tools`、期望最终内容、关键 evidence、允许/禁止动作、正确终止状态和最大预算。这样可以区分“最后抓到了”与“Agent 策略合理”。

另外建立 skill regression cases：同一 observation 输入给固定 Pi/model/toolset 时，验证关键 tool policy 与停止原则；不要求完全相同自然语言或隐藏推理。

再选至少 30 个真实公开 URL 做云端 smoke test，并保存当时的 skill/model/runtime revision 和 evidence。

### 18.4 初始质量门槛

对预先标为可访问且在支持范围内的 fixture，完整正文归档目标 ≥90%；document-like 静态页面中无需 browser 便完成的比例目标 ≥95%；标有可靠 adapter 的 app-like fixture 中，正确选择 adapter 或得到同等/更好证据的比例目标 ≥90%。这是发布门槛，尚非实测结果。

登录墙、验证码、空页面和恶意页面误报“完整归档”的数量必须为 0；安全测试中的跨 library 访问、未经授权发布和密钥外发必须为 0。有限测试集通过不代表所有互联网页面都安全。

所有已提交阅读正文必须能追溯到已保存来源；校验集中的直接引文、数值与代码片段不允许无来源改写。缺失内容必须在 UI 和 manifest 一致显示。

### 18.5 故障注入

覆盖 fetch 后崩溃、R2 写入成功后 D1 失败、模型返回后 checkpoint 前崩溃、浏览器中途关闭、Workflow 创建成功但 outbox 更新失败、向量写入成功但状态更新失败。

覆盖取消期间仍在途、删除后迟到回调、同一 URL 并发收藏、新旧 generation 乱序结束、错误的预算返还、备份期间发生写入、恢复后旧 share token 被访问。

每个测试验证是否丢数据、重复扣预算、产生错误 current pointer、泄露权限或留下不可清理对象；不能只断言请求最终返回 200。

### 18.6 搜索与总结评测

构建至少 30 个中文/英文检索问题，其中包含准确项目名、模糊概念、时间线回忆与原文没有同义表述的查询。留出集目标 Recall@10 ≥90%，同时记录每题结果与失败原因。

总结测试覆盖短文、长文、数字/限定条件、代码、缺失资源和 partial 输入。需要逐条检查关键观点与引用；不能只让同一个模型自评“总结正确”。

### 18.7 核心性能目标

初始性能基准为单个 library 中 10,000 篇文章、约 100,000 个正文块，3 个并发交互请求，另有不超过第 16 节并发上限的后台抓取。测试报告记录实际数据大小、请求位置、冷/热状态与采样窗口。

在这一规模下，保存 API 的 P95 目标 ≤1 秒，前提是只完成持久登记、不等待抓取；关键词搜索 P95 目标 ≤1 秒，语义检索 P95 目标 ≤3 秒。语义检索包括 query Embedding，不包括回答生成；这些仍是待验证目标，受网络、区域和模型影响。

不会为不同网站承诺统一归档完成时间。UI 应持续显示可验证的阶段与失败原因，避免虚构精确完成百分比。


## 19. 里程碑、依赖与任务拆分

**以下任务全部处于未开始状态。** 每个里程碑结束时交付可运行代码、迁移、测试证据与操作说明；不能只完成接口定义便宣称通过。每个任务可直接转成一个 issue，过大的任务再拆子任务。

不建设通用 Agent 平台，但 **Pi runtime + Skills + Tool contracts 是产品核心，不是临时 spike**。本仓库只实现 ReKeep 所需的 Capture/Analysis Agent、skill loader/state adapter、受控 tools 与可恢复执行；不复制 Pi coding-agent、插件市场或任意代码执行环境。

### 19.1 依赖关系

```text
M0 技术与模型验证
        │
        ▼
M1 应用与持久状态基座
        │
        ▼
M2 AI 驱动抓取与可靠归档
        │
        ├──────────────┐
        ▼              ▼
M3 AI 总结与阅读    M4 搜索与问答
        │              │
        └───────┬──────┘
                ▼
M5 分享、导出、删除与恢复
                │
                ▼
M6 安全、质量与发布验收
```

M3 与 M4 可以在 ContentRevision/Chunk 协议稳定后并行。安全约束、预算、幂等和删除标记从 M1/M2 开始实现，M6 是验证与收尾，不是第一次加入这些能力。

### 19.2 M0：验证 Pi + Skills 在 Workers 上的核心路径

**目标**：证明 Pi runtime、skill 驱动 tool loop、Cloudflare Browser Run 与 Workflow recovery 可以在真实 Workers 环境组成我们需要的 CaptureAgent。M0 未通过，不能退回固定 scraper pipeline 冒充 Agent。

**M0-T01 · 建立 staging 验证环境。** 创建隔离的 App/Jobs Worker、D1、R2、AI、Browser 和 Workflow binding，启用受限访问。验收：真实 cloud request 能写 evidence，未认证请求被拒绝。依赖：无。

**M0-T02 · 验证 Pi runtime on Workers。** 最小化运行 `@earendil-works/pi-agent-core`：模型调用、tool calling、event stream、state serialize/restore，并验证 ReKeep 自定义 state/tool adapters。验收：记录 bundle、compat flags、版本 pin 与已知 patch；确认可在 Workers 中稳定恢复 Agent state，且不引入 bash/filesystem/coding tools。依赖：M0-T01。[S32]

**M0-T03 · 建立 capture skill + tool registry 原型。** 把 Stella web skill 的 document/site/render 策略提炼为 `capture/SKILL.md`；实现 skill revision、`web_fetch/site_list/site_run/finish` stub。验收：同一个 Agent loop 能根据不同 fixture 自主选择不同 tool，业务代码中没有 hostname switch 决定路径。依赖：M0-T02。[S31]

**M0-T04 · 验证一次 HTTP + Defuddle 的 web_fetch。** 实现安全 fetch、source preservation、纯 JS extraction 和 observation signals；明确禁止 tool 内 browser/Jina fallback。验收：中文、代码、表格、重定向、截断与异常编码均输出可追溯 source blocks。依赖：M0-T03。

**M0-T05 · 验证 site adapter 模式。** 实现 registry 与至少两个 adapter（一个 GitHub/HN 类结构化来源，一个社交/record 类来源）。验收：Pi 可以先 `site_list` 再 `site_run`，也可以在 adapter 不足时改用 browser/fetch；adapter 结果带 provenance。依赖：M0-T03。

**M0-T06 · 验证 Pi 驱动 Browser episode。** 使用 Browser Run/Playwright，在一个有界 step 内执行 Pi observe → act → observe → capture。验收：至少一个动态正文由 Agent 自主升级 browser 后成功；强制中断后用新 session + 持久 state/evidence 继续，不复用旧 node refs。[S02][S03]

**M0-T07 · 验证 state/recovery 与 revision fencing。** 每 turn 保存 Pi state，固定 runtime/model/skill/toolset revision；在 model 返回、tool 前后和 browser episode 中注入 crash。验收：恢复不重复发布、不悄悄换 skill，旧 generation 不更新 current pointer。依赖：M0-T02 至 M0-T06。

**M0-T08 · 验证存储与检索基础。** 测试 D1 batch、FTS5、中文预分词、Embedding、Vectorize namespace 和异步可见性。验收：library 隔离、索引版本和实际向量维度正确。依赖：M0-T01。

**M0-T09 · 固定 Agent strategy eval。** 建立第 18 节 fixture，记录 first tool、tool sequence、最终质量、预算、browser 增益和错误分类。验收：输出可复现报告，证明 Agent 不只是“每个 URL 都 fetch→browser”。依赖：M0-T03 至 M0-T08。

**退出条件**：真实云端至少跑通 `fetch-only`、`site-adapter-first`、`fetch→browser`、`needs_input` 四类路径；获取策略来自 Pi + skill；Workflow 可恢复；确认无需独立 Go 服务。

### 19.3 M1：应用、数据与可靠任务基座

**目标**：用户可以安全登录、提交收藏，并看到可靠保存的任务记录；即使任务派发失败也不会丢失收藏。

**M1-T01 · 建立工程与持续集成。** 创建第 17 节 monorepo、共享类型、schema 校验、测试、lint 和 staging 部署流程。验收：CI 能构建两个 Worker 和前端，锁定依赖与 compatibility_date。依赖：M0。

**M1-T02 · 认证、API token 与统一授权。** 实现 Access JWT、library 归属、scoped API token、CSRF 和生产旁路保护。验收：负向授权测试通过；任何资源读取和写入都不能只凭客户端 library_id。依赖：M1-T01。

**M1-T03 · 领域迁移与版本规则。** 首批实现 Library、Item、Run、Snapshot、ContentRevision、Job、Outbox、幂等、usage 和 tombstone 表及约束；其余业务表随里程碑增加。验收：重复保存、跨库关联、generation 乱序与相同幂等 key 不同输入均有测试。依赖：M1-T01。

**M1-T04 · 对象写入与 evidence 协议。** 实现 R2 key 生成、条件写入、校验值、私有对象读取和 provenance 协议。验收：对象内容与 hash 一致；模型和客户端不能指定任意 bucket/key；未提交对象不会出现在阅读页。依赖：M1-T02、M1-T03。

**M1-T05 · Outbox、Workflow 派发与补偿。** 实现收藏事务、确定性实例 ID、Cron 补发、状态对账和取消标记。验收：在实例创建前后注入失败，最终只有一条有效业务任务，收藏不丢失。依赖：M1-T03。

**M1-T06 · 配额、资源租约与最小 UI。** 实现 usage 预留/结算、resource_slots、统一维护写屏障，以及收藏表单、列表、任务详情和状态轮询。验收：刷新后状态仍正确；并发任务不会仅靠内存计数；缺少模型配置时显示明确配置错误。依赖：M1-T02 至 M1-T05。

**退出条件**：提交收藏可在 D1 中可靠登记，并由最小测试 Workflow 消费；权限、预算与重试骨架可用。此时属于技术基座，不能称为已完成的产品。

### 19.4 M2：Pi CaptureAgent、Skills 与可信归档

**目标**：把 M0 原型升级为生产级 agent-driven acquisition。新增站点或策略主要通过 skill/adapter 扩展，错误或不完整不能伪装成成功。

**M2-T01 · Capture skill v1 与版本发布机制。** 完成目标、策略原则、source type 判断、quality heuristics、browser/adapter 使用边界、停止条件和安全规则；实现 immutable skill revision。验收：run 可查看当时 skill revision；更新 skill 不改变在途 run。依赖：M1 + M0-T03。

**M2-T02 · 生产级 tool executor 与 Pi state store。** 实现 tool schema、ToolContext、Pi state serialize、bounded observations、events、预算、no-progress、cancel/generation fencing。验收：Agent 进程内存丢失后可从 D1/R2 恢复下一 turn；任何 tool 不能切换 library 或扩大预算。依赖：M1。

**M2-T03 · 生产级 web_fetch / extract / read_source。** 补齐 SSRF、逐跳重定向、timeout、MIME/bytes、Defuddle/source blocks、多 extractor 与 metadata evidence。验收：tool 永远只返回 observation，不自行启动其他获取策略。依赖：M2-T02。

**M2-T04 · Site adapter registry 与首批 adapters。** 建立 descriptor/match/run/provenance contract，加入首批真实来源；实现 `site_list/site_run`。验收：不需要修改 Agent loop 就能增加 adapter；错误 adapter 结果可以被 Agent 拒绝或换策略。依赖：M2-T02。

**M2-T05 · Agent-driven Browser Run tools。** 实现 open/observe/act/capture/close、read-only action policy、episode lifecycle、重新观察与故障恢复。验收：动态展开/分页样本成功；每次动作可以追溯到 Pi event 和 observation；没有 tool 内隐藏的固定点击脚本。依赖：M2-T02。

**M2-T06 · 来源忠实提取与 ArchiveService。** 实现 selected block proposal、确定性 Markdown render、独立 verifier、content/resource quality、manifest 和不可变提交。验收：模型生成文字不进入原文；partial 不覆盖旧 complete；R2 成功/D1 失败可对账恢复。依赖：M2-T03 至 M2-T05。

**M2-T07 · Assets、needs-input 与阅读页。** 实现关键媒体归档、补充正文、重新抓取、取消、快照/Agent trail 查看。验收：登录墙进入 needs_input；用户补充后创建 child run；阅读页不执行原站脚本。依赖：M2-T06。

**M2-T08 · Capture strategy eval 与 skill tuning。** 跑第 18 节固定集，对比 first tool、tool sequence、完整度、误报、browser 增益、turn/cost；只通过发布新的 skill revision 改策略，不往业务代码追加站点分支。验收：达到质量门槛，所有失败都有分类。依赖：M2-T01 至 M2-T07。

**退出条件**：输入 URL 后由 Pi 自主选择 fetch/site/browser/read 等组合，最终得到 complete、partial、needs_input 或明确失败；新增一个 site adapter 不需要修改 CaptureAgent 主循环。

### 19.5 M3：AI 总结与日常阅读管理

**目标**：保存并利用 AI 对内容的理解，同时保留用户自己的上下文与历史总结。

**M3-T01 · Pi AnalysisAgent + analyze skill。** 实现受限 toolset、版本化 `analyze/SKILL.md`、正文 outline/按需读取、长文覆盖计划、关键观点、标签与 coverage。验收：Agent 只能读取固定 Snapshot；长文不静默截断；模型失败不影响原文；同 operation key 不产生重复生效的总结。依赖：M2 的归档协议。

**M3-T02 · 引用与质量校验。** 实现引用锚点、直接引文核对、数字/限定条件验证、一次纠正与失败标记。验收：每个有出处的结论能打开正确 ContentRevision 的段落；无证据结果不冒充原文。依赖：M3-T01。

**M3-T03 · 已有总结导入与版本比较。** 支持粘贴外部 AI 回答、注明来源、关联快照、保留生成历史；重新总结创建新版本。验收：导入内容默认标为用户提供，不伪造出处，不覆盖原文或旧总结。依赖：M3-T01。

**M3-T04 · 笔记、收藏夹、标签与阅读时间线。** 实现私人备注、简单笔记编辑、标签来源区分、阅读状态与时间线。验收：并发编辑检测版本冲突；自动标签不删除用户标签；收藏日期和阅读日期不混淆。依赖：M1、M2。

**M3-T05 · 阅读页与总结交互。** 完成原文/总结/笔记分区、引用跳转、处理状态和重做入口。验收：重抓网页与重做总结是不同操作；partial 输入有醒目标记；未完成语义索引不影响阅读。依赖：M3-T02 至 M3-T04。

**退出条件**：AI 总结是可靠保存、有版本、有证据的独立产物；用户能保留已有总结和个人阅读上下文。定期日报/周报不阻塞 MVP。

### 19.6 M4：全文搜索、语义找回与知识库问答

**目标**：既能按项目名精确查找，也能通过模糊记忆找回原文，并生成可核验的跨文章回答。

**M4-T01 · 文本分块与 FTS 更新。** 实现中文/英文预分词、Chunk、标题/正文/笔记索引、变更 outbox 和索引重建。验收：来源文本不被分词结果污染；索引版本可重建；删除和版本切换不会返回旧生效文本。依赖：M2 ContentRevision、M3 笔记协议。

**M4-T02 · Embedding 与 Vectorize。** 实现有界批次、模型维度校验、library namespace、全局唯一向量 ID、upsert 幂等与可见性确认。验收：不同模型索引不混用；写入成功与可查询状态分开显示。依赖：M4-T01。

**M4-T03 · 混合检索与过滤。** 实现 FTS/向量并行召回、RRF、文章去重、日期与来源过滤、D1 最终授权和 fallback。验收：向量不可用时关键词搜索仍可用；任何返回片段先经过有效快照与权限检查。依赖：M4-T02。

**M4-T04 · 有出处的问答。** 实现问题归一化、证据选择、上下文上限、答案 streaming、引文检查和无证据拒答。验收：不存在的来源不能生成合法引用；默认不会联网补写用户阅读历史。依赖：M4-T03、M3-T02。

**M4-T05 · 搜索 UI 与固定评测。** 完成全文/语义状态提示、可读片段、过滤、引用打开与用户反馈。验收：执行固定检索问题，给出 Recall@10、失败明细和实际延迟，不用演示题替代留出评测。依赖：M4-T03、M4-T04。

**退出条件**：关键词、语义、问答三条路径可独立使用；原文、用户笔记和 AIArtifact 在结果中明确区分；通过质量与删除隔离测试。

### 19.7 M5：分享、删除、导出与恢复

**目标**：归档不仅能存进去，还能安全分享、撤销、删除、导出并在故障后恢复。

**M5-T01 · 版本固定分享。** 实现 token hash、有效期、允许字段、全文 opt-in、独立公开入口和附件授权。验收：分享不会泄露私人笔记；撤销后页面与附件均拒绝访问；原站失效不影响归档分享。依赖：M2、M3、M1 授权。

**M5-T02 · 删除与迟到任务清理。** 实现 tombstone、share 撤销、generation fencing、PurgeWorkflow、向量/资源扫尾。验收：删除后搜索、问答、直读和分享立即拒绝，迟到 Workflow 不能使内容复活。依赖：M4、M5-T01。

**M5-T03 · 单篇与全库导出。** 实现版本化 Markdown/JSONL、资源映射、分卷、manifest 和 hash 校验。验收：导出不依赖访问原站；不包含秘密；超出内存时仍可流式/分批完成。依赖：M2、M3。

**M5-T04 · 一致备份与恢复工具。** 实现写屏障、非 FTS 主表逻辑导出、对象保护引用、独立 backup bucket、恢复与索引重建。验收：60 秒屏障超时正确失败并解锁；隔离恢复校验数据；不默认激活旧 token、share 或在途任务。依赖：M5-T02、M5-T03。

**M5-T05 · 清理、保留与管理界面。** 实现未引用对象 GC、失败尝试保留、备份清理、存储/用量概览、API token 管理与维护提示。验收：已提交快照引用的 run 对象不被过期删除；访问和导出不因 AI 额度耗尽而被阻断。依赖：M5-T01 至 M5-T04。

**退出条件**：完成一次完整导出与隔离恢复演练；分享撤销和删除不依赖异步向量清理速度；没有只有平台运行时才能读取的唯一内容副本。

### 19.8 M6：发布验收与运行手册

**目标**：交付可用于真实个人阅读的版本，并明确支持范围、失败方式和日常运维方法。

**M6-T01 · 端到端回归与真实 URL 验证。** 固定 staging 版本，执行全部 fixture、留出问题和 30 个公开 URL。验收：报告包含样本、时间、模型/代码版本、成功分母、费用和失败原因。依赖：M2 至 M5。

**M6-T02 · 安全与故障演练。** 执行权限、提示注入、网络跳转、原文脚本、任务中断、重复事件、乱序发布和删除后恢复测试。验收：第 18 节安全与故障注入门槛通过。依赖：M6-T01。

**M6-T03 · 预算与性能校准。** 测试约定规模和并发，校准模型上下文、browser 占比、默认预算、费用告警和队列公平性。验收：报告实测值；未达延迟目标时限定规模或优化，不能宣称理论数值已达到。依赖：M6-T01。

**M6-T04 · 部署、升级与事故手册。** 编写从空账户资源配置到部署的步骤、模型替换、Workflow 升级、失败重试、密钥轮换、恢复与回滚手册。验收：按文档完成 staging 重建，旧 Workflow 版本不因部署被破坏。依赖：M5-T04。

**M6-T05 · Release candidate 与私人试用。** 通过第 20 节检查，生成版本与已知限制清单，使用真实个人收藏验证日常流程。验收：阻断问题有明确结论；未支持来源显示可理解状态；不自动开放多用户或匿名抓取入口。依赖：M6-T01 至 M6-T04。

**退出条件**：达到定义的 MVP 完成标准，可以稳定处理主动收藏；仍不宣称“可抓取所有网站”或“AI 永不判断错误”。

### 19.9 MVP 之后的独立工作包

**P1 · 浏览器扩展。** 主动提交当前页面、选择文本和正文提取结果；针对登录态内容重新审查权限与隐私。复用既有输入协议，不改写归档核心。

**P2 · 日报与周报。** 实现 DigestWorkflow 和按用户时区的到期派发，按已读/收藏范围生成有出处的回顾；与核心收藏任务使用不同预算优先级。

**P3 · 新内容来源。** PDF、平台专用 API、X/视频转录、GitHub 等分别定义来源适配器、权限条件、大小预算和质量 fixture，不以“支持 URL”概括所有格式。

**P4 · 团队化与规模化。** 只有明确需求后再加入多人共享知识库、协作权限、Queue 派发、D1 分库或独立搜索分片；不把这些作为个人版发布前置条件。

## 20. MVP 完成标准

**20.1 获取真正由 Agent 驱动。** 对安全允许的输入，Pi CaptureAgent 依据 capture skill、URL 类型与已有 evidence 自主选择 `web_fetch / site adapter / browser / read`；业务代码不存在强制 `fetch→browser` pipeline。Browser 可观察与交互，site adapter 是一等能力，新增来源不需要修改 Agent 主循环。

**20.2 原文可信。** 保存的阅读正文能追溯到来源块，模型生成文字不混入原文。来源时间、版本、质量和缺失资源可查看，已知缺口不标完整。

**20.3 归档独立于源站。** 在阻断原站访问的测试中，仍能读取已保存正文与关键资源，并打开引用锚点和授权分享。未归档附件明确标记外链。

**20.4 AI 从第一版可用。** 自动总结、关键观点、已有总结保存和有出处问答可运行；模型失败、额度耗尽或 Embedding 延迟不影响原文读取与导出。

**20.5 找回能力可核验。** FTS、语义检索和时间/标签过滤可用，达到固定留出集门槛；答案引用固定快照，不把实时网页或模型常识伪装成个人阅读记录。

**20.6 故障不会静默丢内容。** Outbox 可补发、任务可重试、旧 generation 不覆盖新版本、浏览器中断可恢复；重复外部费用有记录，不能只靠唯一键宣称没有重复调用。

**20.7 隐私和撤销生效。** 未授权用户不能访问私库；分享字段受控且可撤销；删除即时在 API 层生效；原网页不能执行应用权限操作。

**20.8 数据可迁移、可恢复。** 至少完成一次全库逻辑导出和隔离恢复，正文/笔记/总结校验通过，索引可重建；恢复不复活旧访问凭据与后台任务。

**20.9 成本与运行边界清楚。** 抓取、browser、模型和资源都有可恢复预算；用户能看到任务失败原因与消耗；单人版不开放匿名爬取入口。

**20.10 交付物完整。** 代码、migration、配置样例、测试报告、来源适配说明、部署/升级/恢复手册和已知限制齐备。本 Plan 中的目标数值不能冒充实际测试成绩。

## 21. 架构决策与需要验证的风险

### ADR-01：全 Cloudflare，但不把所有能力放进一个请求

Workers 承载代码，Workflows 负责异步恢复，D1/R2 保存业务结果。无需独立 Go 服务，但也不能把浏览器与模型长链放进普通 HTTP 请求或仅依赖 waitUntil。

**代价**：需要维护 Workflow 版本、异步状态与平台绑定。通过清晰的数据格式、逻辑模块和导出降低迁移成本，不提前建设跨云兼容层。

### ADR-02：Pi Agent owns orchestration；Skill encodes strategy；Tools expose capability

CaptureAgent 使用 Pi runtime。`capture/SKILL.md` 负责 document/site/browser 的策略知识；tools 只执行一次获取、观察、浏览器动作或归档动作。ArchiveService 保留安全与最终发布权。

**代价**：需要 skill revision、Agent state persistence、tool contract 和 strategy eval。但新增来源/经验不再不断扩张 pipeline 分支。

### ADR-03：Fetch 是 document 的默认策略，不是系统硬约束；Site Adapter 是一等路径

普通 document 通常从 `web_fetch` 开始，但 app-like records 可以直接选择 `site_run`；用户提交正文或明确动态页面也无需为了形式先 fetch。`web_fetch` 内部不自动 fallback 到 browser/Jina。

**代价**：Agent 可能偶尔选错首个工具，需要多一次 turn；通过 skill/eval/预算优化，而不是重新写 hostname switch。

### ADR-04：Browser actions 由 Pi 决策，但 session lifecycle 在一个 durable episode 内管理

逻辑上 browser 是 open/observe/act/capture/close tools；物理上 page handle 只在有界 Workflow step 内存活。step 内运行 Pi browser loop，崩溃后新建 session 并从持久 evidence 重新观察。

**代价**：中断可能重复模型和浏览器成本，但不会伪造可序列化的 page 状态，也不会把 browser planner 隐藏进 tool。

### ADR-05：D1 + R2 + Vectorize 各有单一职责

D1 管有效状态与权限，R2 保存可导出的内容，Vectorize 是可重建的辅助索引。一次归档通过 manifest 和 D1 发布，接受后台清理未引用对象。

**代价**：没有跨服务事务，需要 outbox、对账、tombstone 和 GC。不能把“一个 API 成功”当成三个系统同时提交。

### ADR-06：Pi/runtime 与 Skills 保持平台中立，Cloudflare 通过 adapters 接入

Capture/Analysis Agent、skills 和 tool contracts 不直接依赖 D1/R2/Browser Run。MVP 只实现 Cloudflare adapters，但这允许未来把同一 Agent core 运行在 celld，替换存储、browser 和 workflow adapter。

**代价**：会多一层小型 interface，但不能演变成“大一统跨云框架”；只有已经存在的产品能力才抽象。

### ADR-07：Skill revision 是生产代码的一部分

Skill 改动可以改变 Agent 获取策略，因此必须像代码一样 version、review、eval、pin 到 run。不能在线直接编辑一段 prompt 并让所有在途任务立即改变行为。

**代价**：发布 skill 需要 regression eval；收益是事故可追溯、行为可重放、策略迭代无需改 tool implementation。

### 21.1 必须在 M0/M2 关闭的风险

**Pi Worker 兼容性**：`@earendil-works/pi-agent-core` 与所选 Pi provider 在 Cloudflare Workers bundle/runtime 是否稳定。M0 固定版本、compat flags、state serialization/recovery 行为和兼容性证据。

**模型工具调用质量**：候选模型是否能理解 capture skill、选择合理首个 tool、在 observation 不足时切换策略并克制 browser 升级。通过 strategy 留出集选定配置，不依赖模型名称或宣传指标。

**Skill 漂移与过拟合**：skill tuning 可能修好少数站点却破坏通用行为。每个 revision 跑固定 + 留出集，并记录 tool sequence/cost，而不仅看最终有没有抓到。

**网络安全与浏览器策略**：实际运行时能否实施所需 URL/重定向/资源限制。无法验证的危险路径默认拒绝，不给模型“绕过限制”的工具。

**浏览器生命周期**：使用的 Playwright 版本、close 行为、会话超时、异常退出计费与清理必须云端测试。M0 记录真实证据。

**全文提取忠实度**：parser 对代码、表格、中文、分页和虚拟滚动的处理是否可追溯。无法完整获得的结构标 partial，不能由模型补全。

**一致恢复**：FTS 存在时的逻辑备份、短写屏障和跨 R2 清单是否正确。若库规模超过初始屏障能力，备份失败要可见，并升级方案后再扩大使用。

**费用偏差**：浏览器内等待模型、重试、未知用量和长文分析可能显著超过单次摘要费用。使用应用计量与 Cloudflare 账单对账，不声称应用限额能封顶全账户费用。

### 21.2 当前已确定，不再作为实现前置问题

第一版以个人私有知识库为主；主动保存 URL 和文本；生产后端全 Cloudflare；Agent runtime 使用 Pi；获取策略由 versioned capture skill 驱动；`web_fetch / site adapter / browser` 是并列 capabilities，不存在硬编码 fetch→browser pipeline；原文、提取版本和 AIArtifact 分离；分享可撤销；全文导出与恢复属于 MVP。

产品品牌、视觉风格、最终模型优选和后续来源适配可独立决定。它们不阻碍 M0 验证和 M1 基座开发。

**实施起点：先执行 M0-T01 至 M0-T09，用真实云端结果确认 Pi runtime、skills、site adapters、Browser Run、状态恢复与安全边界，再进入 M1 和 M2。**

## 22. 参考来源

以下均为本次设计查阅的官方文档或项目原始资料。平台能力、限制与价格可能变化，实施时应以固定版本和当前账号配置再次核对。引用用于支持平台事实；本文任务、预算和质量门槛是本项目的设计决定。

**[S01] · Workers Static Assets。** 前端静态资源与 Worker 一体部署。

**[S02] · Browser Run 概览。** 浏览器会话、Quick Actions、Agent 使用方式与产品边界。

**[S03] · Browser Run Playwright。** Cloudflare Playwright 集成、会话连接与关闭语义。

**[S04] · Browser Run FAQ。** 机器人识别、站点兼容性和运行限制。

**[S05] · Browser Run Stagehand。** 集成方式与支持版本，不能直接假设兼容最新版。

**[S06] · Workflows Durable Agents。** 将模型决策与工具调用作为可持久化步骤的官方模式。

**[S07] · Rules of Workflows。** 步骤、副作用、重试与幂等的约束。

**[S08] · Workflows Limits。** 步骤输出、状态保留、CPU 与执行限制。

**[S09] · D1 SQL Statements。** SQLite 语法与 FTS5 支持。

**[S10] · D1 Database API。** batch 批量事务与运行时 API。

**[S11] · Browser Run Limits。** 浏览器空闲生命周期与资源限制。

**[S12] · R2 Workers API。** 对象读写、条件写入与完整性校验。

**[S13] · D1 Limits。** 单库容量、行大小、绑定参数和 Time Travel 保留。

**[S14] · Vectorize Query Vectors。** 向量查询、namespace 与 metadata 过滤。

**[S15] · Vectorize Insert Vectors。** 插入、upsert 与异步索引可见性。

**[S16] · Vectorize Limits。** topK、ID 长度、维度和批量操作上限。

**[S17] · Workers AI Qwen3-30B-A3B-FP8。** 本文候选模型的上下文和 function calling 支持。

**[S18] · Workers AI BGE-M3。** Embedding 模型接口；1024 维设计同时参考 [BAAI 原始模型卡](https://huggingface.co/BAAI/bge-m3)。实际 binding 输出仍需验证。

**[S19] · Cloudflare Access JWT Validation。** Worker 应用校验 Access JWT 的方法。

**[S20] · OWASP SSRF Prevention。** 请求目标、重定向与网络访问风险。

**[S21] · OWASP LLM Prompt Injection Prevention。** 不可信网页内容与工具权限边界。

**[S22] · AI Gateway Logging。** 请求/响应日志、payload 与用量记录控制。

**[S23] · Workers Limits。** isolate 内存、请求与执行限制。

**[S24] · SQLite FTS5。** tokenizer、全文索引与查询语义。

**[S25] · D1 Time Travel。** 数据库时间点恢复能力与边界。

**[S26] · Browser Run Pricing。** 浏览器会话时间、并发与 Quick Actions 计费区别。

**[S27] · Workflows Pricing。** CPU、步骤与状态存储计费。

**[S28] · Workers AI Data Usage。** 模型服务的数据使用说明。

**[S29] · D1 Import / Export。** 导出行为、virtual table 限制与迁移注意事项。

**[S30] · Workers HTMLRewriter。** Workers 内流式 HTML 处理 API。

**[S31] · CherryHQ Stella Web Skill。** 本项目借鉴其 document fetch、site records、rendered page 与 site script 的策略划分；ReKeep 不照搬其 fetch 内部固定 fallback 和 Bun/Lightpanda 本地执行实现。

**[S32] · Pi Agent Core。** `@earendil-works/pi-agent-core` 的 stateful agent、tool execution、state management，以及 runtime-specific storage backend 设计。


[S01]: https://developers.cloudflare.com/workers/static-assets/
[S02]: https://developers.cloudflare.com/browser-run/
[S03]: https://developers.cloudflare.com/browser-run/playwright/
[S04]: https://developers.cloudflare.com/browser-run/faq/
[S05]: https://developers.cloudflare.com/browser-run/stagehand/
[S06]: https://developers.cloudflare.com/workflows/get-started/durable-agents/
[S07]: https://developers.cloudflare.com/workflows/build/rules-of-workflows/
[S08]: https://developers.cloudflare.com/workflows/reference/limits/
[S09]: https://developers.cloudflare.com/d1/sql-api/sql-statements/
[S10]: https://developers.cloudflare.com/d1/worker-api/d1-database/
[S11]: https://developers.cloudflare.com/browser-run/limits/
[S12]: https://developers.cloudflare.com/r2/api/workers/workers-api-reference/
[S13]: https://developers.cloudflare.com/d1/platform/limits/
[S14]: https://developers.cloudflare.com/vectorize/best-practices/query-vectors/
[S15]: https://developers.cloudflare.com/vectorize/best-practices/insert-vectors/
[S16]: https://developers.cloudflare.com/vectorize/platform/limits/
[S17]: https://developers.cloudflare.com/workers-ai/models/qwen3-30b-a3b-fp8/
[S18]: https://developers.cloudflare.com/workers-ai/models/bge-m3/
[S19]: https://developers.cloudflare.com/cloudflare-one/access-controls/applications/http-apps/authorization-cookie/validating-json/
[S20]: https://cheatsheetseries.owasp.org/cheatsheets/Server_Side_Request_Forgery_Prevention_Cheat_Sheet.html
[S21]: https://cheatsheetseries.owasp.org/cheatsheets/LLM_Prompt_Injection_Prevention_Cheat_Sheet.html
[S22]: https://developers.cloudflare.com/ai-gateway/observability/logging/
[S23]: https://developers.cloudflare.com/workers/platform/limits/
[S24]: https://www.sqlite.org/fts5.html
[S25]: https://developers.cloudflare.com/d1/reference/time-travel/
[S26]: https://developers.cloudflare.com/browser-run/pricing/
[S27]: https://developers.cloudflare.com/workflows/reference/pricing/
[S28]: https://developers.cloudflare.com/workers-ai/platform/data-usage/
[S29]: https://developers.cloudflare.com/d1/best-practices/import-export-data/
[S30]: https://developers.cloudflare.com/workers/runtime-apis/html-rewriter/
[S31]: https://github.com/CherryHQ/stella/blob/main/plugins/agent/web/skills/web/SKILL.md
[S32]: https://github.com/earendil-works/pi/blob/main/packages/agent/README.md
