import { Icon } from "@iconify/react";
import { ActionIcon, CopyButton, Tooltip } from "@mantine/core";

export function CopyBtn({ data, ...props }) {
  return (
    <CopyButton value={data} timeout={2000} {...props}>
      {({ copied, copy }) => (
        <Tooltip label={copied ? "Copied" : "Copy"} withArrow position="right">
          <ActionIcon
            color={copied ? "teal" : "gray"}
            variant="subtle"
            onClick={copy}
          >
            {copied ? (
              <Icon icon="tabler:check" />
            ) : (
              <Icon icon="tabler:copy" />
            )}
          </ActionIcon>
        </Tooltip>
      )}
    </CopyButton>
  );
}
