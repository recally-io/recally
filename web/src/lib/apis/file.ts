import useSWRMutation from "swr/mutation";
import fetcher from "./fetcher";

export interface GetPresignedURLsRequest {
	fileName: string;
	fileType: string;
	action: "PUT" | "GET";
	expiration?: number;
}

export interface GetPresignedURLsResponse {
	presigned_url: string;
	object_key: string;
	public_url: string;
}

export interface FileError {
	message: string;
}

// Get presigned URLs for file upload/download
export const useGetPresignedURLs = () => {
	return useSWRMutation<
		GetPresignedURLsResponse,
		FileError,
		string,
		GetPresignedURLsRequest
	>("/api/v1/files/presigned-urls", async (url, { arg }) => {
		const queryParams = new URLSearchParams({
			file_name: arg.fileName,
			file_type: arg.fileType,
			action: arg.action,
			expiration: arg.expiration?.toString() || "3600",
		}).toString();

		return await fetcher(`${url}?${queryParams}`);
	});
};

// Delete file by ID
export const useDeleteFile = () => {
	return useSWRMutation<void, FileError, string>(
		"/api/files",
		async (url, { arg: fileId }) => {
			return await fetcher(`${url}/${fileId}`);
		},
	);
};
