import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient();

/**
 * Makes an HTTP request to the specified path with the given method, parameters, and body.
 * @param {string} path - The path to send the request to.
 * @param {string} [method="GET"] - The HTTP method to use for the request. Defaults to "GET".
 * @param {Object} [params={}] - The query parameters to include in the request.
 * @param {Object} [body={}] - The request body to send.
 * @returns {Promise<Object>} - A promise that resolves to the JSON response from the server.
 * @throws {Error} - If the request fails or the response status is not OK.
 */
export async function request(method, path, params = null, body = null) {
  const options = {
    method: method,
    headers: {
      "Content-Type": "application/json",
    },
  };

  if (params !== null && params !== undefined) {
    path += "?" + new URLSearchParams(params).toString();
  }

  if (body !== null && body !== undefined) {
    options.body = JSON.stringify(body);
  }

  const res = await fetch(path, options);
  if (!res.ok) {
    throw new Error(
      `Error ${method} ${path}. Response status: ${res.status} ${res.statusText}`,
    );
  }
  if (res.status === 204) {
    return {};
  }
  return await res.json();
}

/**
 * Sends a GET request to the specified path with optional parameters.
 * @param {string} path - The path to send the GET request to.
 * @param {object} params - Optional parameters to include in the request.
 * @returns {Promise<Object>} - A promise that resolves to the JSON response from the server.
 * @throws {Error} - If the request fails or the response status is not OK.
 */
export async function get(path, params = null) {
  return await request("GET", path, params);
}

/**
 * Sends a POST request to the specified path with the given body.
 * @param {string} path - The path to send the POST request to.
 * @param {Object} [params=null] - The query parameters to include in the request.
 * @param {object} [body=null] - The request body to send.
 * @returns {Promise<Object>} - A promise that resolves to the JSON response from the server.
 * @throws {Error} - If the request fails or the response status is not OK.
 */
export async function post(path, params = null, body = null) {
  return await request("POST", path, params, body);
}

/**
 * Sends a PUT request to the specified path with the given body.
 * @param {string} path - The path to send the POST request to.
 * @param {Object} [params=null] - The query parameters to include in the request.
 * @param {object} [body=null] - The request body to send.
 * @returns {Promise<Object>} - A promise that resolves to the JSON response from the server.
 * @throws {Error} - If the request fails or the response status is not OK.
 */
export async function put(path, params = null, body = null) {
  return await request("PUT", path, params, body);
}

/**
 * Sends a DELETE request to the specified path.
 * @param {string} path - The path to send the DELETE request to.
 * @param {Object} [params=null] - The query parameters to include in the request.
 * @returns {Promise<Object>} - A promise that resolves to the JSON response from the server.
 * @throws {Error} - If the request fails or the response status is not OK.
 */
export async function del(path, params = null) {
  return await request("DELETE", path, params);
}

export const listToolsKey = ["list-assistant-tools"];
export async function listTools() {
  const res = await get("/api/v1/assistants/tools");
  let data = res.data || [];
  data = data.map((tool) => tool.name);
  return data;
}

export const listModelsKey = ["list-assistant-models"];
export async function listModels() {
  const res = await get("/api/v1/assistants/models");
  return res.data || [];
}

/**
 * Retrieves a presigned URL for file operations.
 * @param {Object} options - The options for getting the presigned URL.
 * @param {string} options.assistantId - The ID of the assistant.
 * @param {string} options.threadId - The ID of the thread.
 * @param {string} options.fileName - The name of the file.
 * @param {string} options.fileType - The type of the file.
 * @param {string} options.action - The action to be performed (e.g., 'read', 'write').
 * @param {number} options.expiration - The expiration time for the presigned URL.
 * @returns {Promise<Object>} - A promise that resolves to an object containing the public URL and presigned URL.
 */
export async function getPresignedUrl({
  assistantId,
  threadId,
  fileName,
  fileType,
  action = "PUT",
  expiration = 3600,
}) {
  const params = new URLSearchParams({
    assistant_id: assistantId,
    thread_id: threadId,
    file_name: fileName,
    file_type: fileType,
    action: action,
    expiration: expiration,
  });
  const res = await get(`/api/v1/files/presigned-urls`, params);
  return {
    publicUrl: res.data.public_url,
    preSignedURL: res.data.presigned_url,
  };
}

/**
 * Uploads a file using a presigned URL.
 * @param {Object} options - The options for uploading the file.
 * @param {string} options.preSignedURL - The presigned URL for uploading.
 * @param {File} options.file - The file to be uploaded.
 * @param {string} options.publicUrl - The public URL of the file after upload.
 * @returns {Promise<string>} - A promise that resolves to the public URL of the uploaded file.
 * @throws {Error} If the file upload fails.
 */
export async function uploadFile({ preSignedURL, file, publicUrl }) {
  const response = await fetch(preSignedURL, {
    method: "PUT",
    body: file,
    headers: { "Content-Type": file.type },
  });
  if (!response.ok) throw new Error("Failed to upload file");
  return publicUrl;
}

/** * Posts an attachment to an assistant or a thread.
 * @async
 * @param {Object} options - The options for posting the attachment.
 * @param {string} options.assistantId - The ID of the assistant.
 * @param {string} [options.threadId] - The ID of the thread (optional).
 * @param {string} options.type - The type of the attachment.
 * @param {string} options.name - The name of the attachment.
 * @param {string} options.publicUrl - The public URL of the attachment.
 * @param {Object} options.docs - Additional documentation for the attachment.
 * @returns {Promise<Object>} The data returned from the server after posting the attachment.
 */
export async function postAttachment({
  assistantId,
  threadId,
  type,
  name,
  publicUrl,
  docs,
}) {
  let url = "";
  if (threadId) {
    url = `/api/v1/assistants/${assistantId}/threads/${threadId}/attachments`;
  } else {
    url = `/api/v1/assistants/${assistantId}/attachments`;
  }

  const res = await post(url, null, {
    type: type,
    name: name,
    url: publicUrl,
    docs: docs,
  });
  return res.data;
}
