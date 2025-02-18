import { Button } from "@/components/ui/button";
import {
	ChatBubble,
	ChatBubbleAvatar,
	ChatBubbleMessage,
} from "@/components/ui/chat-bubble";
import { ChatInput } from "@/components/ui/chat-input";
import { ChatMessageList } from "@/components/ui/chat-message-list";
import {
	ExpandableChat,
	ExpandableChatBody,
	ExpandableChatFooter,
	ExpandableChatHeader,
} from "@/components/ui/expandable-chat";
import { Bot, CornerDownLeft } from "lucide-react";
import { FormEvent, useState } from "react";

export function ExpandableChatDemo() {
	const [messages, setMessages] = useState([
		{
			id: 1,
			content: "Hello! How can I help you today?",
			sender: "ai",
		},
		{
			id: 2,
			content: "I have a question about the component library.",
			sender: "user",
		},
		{
			id: 3,
			content: "Sure! I'd be happy to help. What would you like to know?",
			sender: "ai",
		},
	]);

	const [input, setInput] = useState("");
	const [isLoading, setIsLoading] = useState(false);

	const handleSubmit = (e: FormEvent) => {
		e.preventDefault();
		if (!input.trim()) return;

		setMessages((prev) => [
			...prev,
			{
				id: prev.length + 1,
				content: input,
				sender: "user",
			},
		]);
		setInput("");
		setIsLoading(true);

		setTimeout(() => {
			setMessages((prev) => [
				...prev,
				{
					id: prev.length + 1,
					content: "This is an AI response to your message.",
					sender: "ai",
				},
			]);
			setIsLoading(false);
		}, 1000);
	};

	// const handleAttachFile = () => {
	// 	//
	// };

	// const handleMicrophoneClick = () => {
	// 	//
	// };

	return (
		<div className="h-[600px] relative">
			<ExpandableChat
				size="lg"
				position="bottom-right"
				icon={<Bot className="h-6 w-6" />}
			>
				<ExpandableChatHeader className="flex-col text-center justify-center">
					<h1 className="text-xl font-semibold">Chat with AI ✨</h1>
					<p className="text-sm text-muted-foreground">
						Ask me anything about the components
					</p>
				</ExpandableChatHeader>

				<ExpandableChatBody>
					<ChatMessageList>
						{messages.map((message) => (
							<ChatBubble
								key={message.id}
								variant={message.sender === "user" ? "sent" : "received"}
							>
								<ChatBubbleAvatar
									className="h-8 w-8 shrink-0"
									src={
										message.sender === "user"
											? "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=64&h=64&q=80&crop=faces&fit=crop"
											: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=64&h=64&q=80&crop=faces&fit=crop"
									}
									fallback={message.sender === "user" ? "US" : "AI"}
								/>
								<ChatBubbleMessage
									variant={message.sender === "user" ? "sent" : "received"}
								>
									{message.content}
								</ChatBubbleMessage>
							</ChatBubble>
						))}

						{isLoading && (
							<ChatBubble variant="received">
								<ChatBubbleAvatar
									className="h-8 w-8 shrink-0"
									src="https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=64&h=64&q=80&crop=faces&fit=crop"
									fallback="AI"
								/>
								<ChatBubbleMessage isLoading />
							</ChatBubble>
						)}
					</ChatMessageList>
				</ExpandableChatBody>

				<ExpandableChatFooter>
					<form
						onSubmit={handleSubmit}
						className="relative rounded-lg border bg-background focus-within:ring-1 focus-within:ring-ring p-1"
					>
						<ChatInput
							value={input}
							onChange={(e) => setInput(e.target.value)}
							placeholder="Type your message..."
							className="min-h-12 resize-none rounded-lg bg-background border-0 p-3 shadow-none focus-visible:ring-0"
						/>
						<div className="flex items-center p-3 pt-0 justify-between">
							{/* <div className="flex">
								<Button
									variant="ghost"
									size="icon"
									type="button"
									onClick={handleAttachFile}
								>
									<Paperclip className="size-4" />
								</Button>

								<Button
									variant="ghost"
									size="icon"
									type="button"
									onClick={handleMicrophoneClick}
								>
									<Mic className="size-4" />
								</Button>
							</div> */}
							<Button type="submit" size="sm" className="ml-auto gap-1.5">
								Send Message
								<CornerDownLeft className="size-3.5" />
							</Button>
						</div>
					</form>
				</ExpandableChatFooter>
			</ExpandableChat>
		</div>
	);
}
