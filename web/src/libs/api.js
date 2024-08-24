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
