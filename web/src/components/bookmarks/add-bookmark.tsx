import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useToast } from "@/hooks/use-toast";
import { useBookmarkMutations } from "@/lib/apis/bookmarks";
import {
	type GetPresignedURLsRequest,
	type GetPresignedURLsResponse,
	useGetPresignedURLs,
} from "@/lib/apis/file"; // Import the useGetPresignedURLs hook
import { cn } from "@/lib/utils";
import { Link2, PlusCircle, Upload, X } from "lucide-react";
import { useRef, useState } from "react";

export default function AddBookmarkModal() {
	const { trigger: getPresignedUrls } = useGetPresignedURLs();
	const { createBookmark } = useBookmarkMutations();
	const { toast } = useToast();
	const [activeTab, setActiveTab] = useState<"url" | "file">("url");
	const [open, setOpen] = useState(false);
	const [url, setUrl] = useState("");
	const [urlError, setUrlError] = useState("");
	const [file, setFile] = useState<File | null>(null);
	const [isLoading, setIsLoading] = useState(false);
	const fileInputRef = useRef<HTMLInputElement>(null);
	const [dragActive, setDragActive] = useState(false);

	const validateUrl = (url: string) => {
		try {
			new URL(url);
			setUrlError("");
			return true;
		} catch {
			setUrlError("Please enter a valid URL");
			return false;
		}
	};

	const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
		if (e.target?.files?.[0]) {
			const selectedFile = e.target.files[0];
			if (selectedFile.size > 50 * 1024 * 1024) {
				// 50MB limit
				toast({
					title: "File too large",
					description: "Please select a file smaller than 50MB",
					variant: "destructive",
				});
				return;
			}
			setFile(selectedFile);
		}
	};

	const handleUploadFile = async (
		selectedFile: File,
	): Promise<GetPresignedURLsResponse> => {
		try {
			console.log(file);
			const params: GetPresignedURLsRequest = {
				fileName: selectedFile.name,
				fileType: selectedFile.type,
				action: "PUT",
				expiration: 3600,
			};
			// Use the hook's mutate function to fetch presigned URLs
			const data = await getPresignedUrls(params);
			if (!data) {
				throw new Error("Failed to get presigned URL");
			}
			const uploadRes = await fetch(data.presigned_url, {
				method: "PUT",
				headers: {
					"Content-Type": selectedFile.type,
				},
				body: selectedFile,
			});
			if (!uploadRes.ok) {
				throw new Error("Upload failed");
			}
			toast({
				title: "File uploaded successfully",
				description: "Your file has been uploaded.",
				variant: "default",
			});
			return data;
		} catch (error: any) {
			toast({
				title: "Upload error",
				description: error.message || "Failed to upload file",
				variant: "destructive",
			});
			throw error;
		}
	};

	const handleDrag = (e: React.DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		if (e.type === "dragenter" || e.type === "dragover") {
			setDragActive(true);
		} else if (e.type === "dragleave") {
			setDragActive(false);
		}
	};

	const handleDrop = (e: React.DragEvent) => {
		e.preventDefault();
		e.stopPropagation();
		setDragActive(false);

		if (e.dataTransfer.files?.[0]) {
			const droppedFile = e.dataTransfer.files[0];
			if (droppedFile.size > 50 * 1024 * 1024) {
				toast({
					title: "File too large",
					description: "Please select a file smaller than 50MB",
					variant: "destructive",
				});
				return;
			}
			setFile(droppedFile);
		}
	};

	const handleSubmitCreateBookmark = async (e: React.FormEvent) => {
		e.preventDefault();
		if (activeTab === "url" && (!url || !validateUrl(url))) return;
		if (activeTab === "file" && !file) return;

		setIsLoading(true);
		try {
			if (activeTab === "url") {
				await createBookmark({ url });
				setUrl("");
			} else {
				console.log("start upload file");
				const resp = await handleUploadFile(file!);
				console.log(resp);
				// Get file extension and validate it
				const extension = file?.name?.split(".").pop()?.toLowerCase();
				const fileType = file?.type || `application/${extension}`;

				await createBookmark({
					url: resp.public_url,
					title: file?.name || "Untitled",
					type: extension,
					s3_key: resp.object_key,
					metadata: {
						file: {
							name: file?.name || "Untitled",
							extension: extension || "unknown",
							mime_type: fileType,
							size: file?.size || 0,
						},
					},
				});
				setFile(null);
				if (fileInputRef.current) {
					fileInputRef.current.value = "";
				}
			}
			toast({
				title: "Success",
				description: "Bookmark created successfully",
			});
			setOpen(false);
		} catch (error) {
			console.log(error);
			toast({
				title: "Error",
				description: "Failed to create bookmark. Please try again.",
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
		}
	};

	const resetForm = () => {
		setUrl("");
		setUrlError("");
		setFile(null);
		if (fileInputRef.current) {
			fileInputRef.current.value = "";
		}
	};

	return (
		<Dialog
			open={open}
			onOpenChange={(newOpen) => {
				setOpen(newOpen);
				if (!newOpen) resetForm();
			}}
		>
			<DialogTrigger asChild>
				<Button variant="ghost" size="icon" className="h-7 w-7">
					<PlusCircle className="size-6" />
				</Button>
			</DialogTrigger>
			<DialogContent className="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>Add New Bookmark</DialogTitle>
					<DialogDescription>
						Enter the details of the new bookmark
					</DialogDescription>
				</DialogHeader>
				<Tabs
					value={activeTab}
					onValueChange={(v) => setActiveTab(v as "url" | "file")}
				>
					<TabsList className="grid w-full grid-cols-2">
						<TabsTrigger value="url" className="flex items-center gap-2">
							<Link2 className="size-4" />
							URL
						</TabsTrigger>
						<TabsTrigger value="file" className="flex items-center gap-2">
							<Upload className="size-4" />
							File
						</TabsTrigger>
					</TabsList>
					<form
						onSubmit={handleSubmitCreateBookmark}
						className="mt-4 space-y-4"
					>
						<TabsContent value="url">
							<div className="space-y-2">
								<Label htmlFor="url">URL</Label>
								<Input
									id="url"
									placeholder="https://example.com"
									value={url}
									onChange={(e) => {
										setUrl(e.target.value);
										if (urlError) validateUrl(e.target.value);
									}}
									className={cn(urlError && "border-red-500")}
								/>
								{urlError && <p className="text-sm text-red-500">{urlError}</p>}
								<p className="text-sm text-muted-foreground">
									Enter the URL of the bookmark
								</p>
							</div>
						</TabsContent>
						<TabsContent value="file">
							<div
								className={cn(
									"space-y-2 rounded-lg border-2 border-dashed p-4 transition-colors",
									dragActive && "border-primary bg-primary/5",
									!dragActive && "border-muted",
								)}
								onDragEnter={handleDrag}
								onDragLeave={handleDrag}
								onDragOver={handleDrag}
								onDrop={handleDrop}
							>
								<Label htmlFor="file">File</Label>
								<Input
									id="file"
									ref={fileInputRef}
									type="file"
									onChange={handleFileChange}
									accept=".pdf,.epub,image/*"
									className="hidden"
								/>
								<div className="flex flex-col items-center justify-center gap-2 py-4">
									{file ? (
										<div className="flex items-center gap-2">
											<span className="text-sm font-medium">{file.name}</span>
											<Button
												type="button"
												variant="ghost"
												size="icon"
												className="h-6 w-6"
												onClick={() => {
													setFile(null);
													if (fileInputRef.current) {
														fileInputRef.current.value = "";
													}
												}}
											>
												<X className="size-4" />
											</Button>
										</div>
									) : (
										<div className="flex flex-col items-center gap-4">
											<Button
												type="button"
												variant="ghost"
												size="lg"
												className="h-24 w-full max-w-xs flex-col gap-2"
												onClick={() => fileInputRef.current?.click()}
											>
												<Upload className="size-10 text-primary" />
												<span className="text-sm font-medium">Upload File</span>
											</Button>
											<p className="text-center text-sm text-muted-foreground">
												Click to upload or drag and drop your file here
											</p>
										</div>
									)}
								</div>
								<p className="text-sm text-muted-foreground">
									Supported formats: PDF, EPUB, and images (max 50MB)
								</p>
							</div>
						</TabsContent>
						<div className="flex justify-end gap-2">
							<DialogClose asChild>
								<Button type="button" variant="outline">
									Cancel
								</Button>
							</DialogClose>
							<Button
								type="submit"
								disabled={
									isLoading ||
									(activeTab === "url" && !url) ||
									(activeTab === "file" && !file)
								}
							>
								{isLoading ? "Adding..." : "Add Bookmark"}
							</Button>
						</div>
					</form>
				</Tabs>
			</DialogContent>
		</Dialog>
	);
}
