/* tslint:disable */
/* eslint-disable */
/**
 * Vibrain API
 * This is a simple API for Vibrain project.
 *
 * The version of the OpenAPI document: 1.0
 * Contact: vibrain@vaayne.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import * as runtime from "../runtime";
import type {
  WebReaderGet200Response,
  WebReaderGet400Response,
  WebSearchGet200Response,
  WebSummaryGet200Response,
} from "../models/index";
import {
  WebReaderGet200ResponseFromJSON,
  WebReaderGet200ResponseToJSON,
  WebReaderGet400ResponseFromJSON,
  WebReaderGet400ResponseToJSON,
  WebSearchGet200ResponseFromJSON,
  WebSearchGet200ResponseToJSON,
  WebSummaryGet200ResponseFromJSON,
  WebSummaryGet200ResponseToJSON,
} from "../models/index";

export interface WebReaderGetRequest {
  url: string;
}

export interface WebReaderPostRequest {
  url: string;
}

export interface WebSearchGetRequest {
  query: string;
}

export interface WebSearchPostRequest {
  query: string;
}

export interface WebSummaryGetRequest {
  url: string;
}

export interface WebSummaryPostRequest {
  url: string;
}

/**
 *
 */
export class ToolsApi extends runtime.BaseAPI {
  /**
   * Read the content of a web page
   * Read web content
   */
  async webReaderGetRaw(
    requestParameters: WebReaderGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<WebReaderGet200Response>> {
    if (requestParameters["url"] == null) {
      throw new runtime.RequiredError(
        "url",
        'Required parameter "url" was null or undefined when calling webReaderGet().',
      );
    }

    const queryParameters: any = {};

    if (requestParameters["url"] != null) {
      queryParameters["url"] = requestParameters["url"];
    }

    const headerParameters: runtime.HTTPHeaders = {};

    const response = await this.request(
      {
        path: `/web/reader`,
        method: "GET",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      WebReaderGet200ResponseFromJSON(jsonValue),
    );
  }

  /**
   * Read the content of a web page
   * Read web content
   */
  async webReaderGet(
    requestParameters: WebReaderGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<WebReaderGet200Response> {
    const response = await this.webReaderGetRaw(
      requestParameters,
      initOverrides,
    );
    return await response.value();
  }

  /**
   * Read the content of a web page
   * Read web content
   */
  async webReaderPostRaw(
    requestParameters: WebReaderPostRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<WebReaderGet200Response>> {
    if (requestParameters["url"] == null) {
      throw new runtime.RequiredError(
        "url",
        'Required parameter "url" was null or undefined when calling webReaderPost().',
      );
    }

    const queryParameters: any = {};

    if (requestParameters["url"] != null) {
      queryParameters["url"] = requestParameters["url"];
    }

    const headerParameters: runtime.HTTPHeaders = {};

    const response = await this.request(
      {
        path: `/web/reader`,
        method: "POST",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      WebReaderGet200ResponseFromJSON(jsonValue),
    );
  }

  /**
   * Read the content of a web page
   * Read web content
   */
  async webReaderPost(
    requestParameters: WebReaderPostRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<WebReaderGet200Response> {
    const response = await this.webReaderPostRaw(
      requestParameters,
      initOverrides,
    );
    return await response.value();
  }

  /**
   * Search the content of a web page
   * Search web content
   */
  async webSearchGetRaw(
    requestParameters: WebSearchGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<WebSearchGet200Response>> {
    if (requestParameters["query"] == null) {
      throw new runtime.RequiredError(
        "query",
        'Required parameter "query" was null or undefined when calling webSearchGet().',
      );
    }

    const queryParameters: any = {};

    if (requestParameters["query"] != null) {
      queryParameters["query"] = requestParameters["query"];
    }

    const headerParameters: runtime.HTTPHeaders = {};

    const response = await this.request(
      {
        path: `/web/search`,
        method: "GET",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      WebSearchGet200ResponseFromJSON(jsonValue),
    );
  }

  /**
   * Search the content of a web page
   * Search web content
   */
  async webSearchGet(
    requestParameters: WebSearchGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<WebSearchGet200Response> {
    const response = await this.webSearchGetRaw(
      requestParameters,
      initOverrides,
    );
    return await response.value();
  }

  /**
   * Search the content of a web page
   * Search web content
   */
  async webSearchPostRaw(
    requestParameters: WebSearchPostRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<WebSearchGet200Response>> {
    if (requestParameters["query"] == null) {
      throw new runtime.RequiredError(
        "query",
        'Required parameter "query" was null or undefined when calling webSearchPost().',
      );
    }

    const queryParameters: any = {};

    if (requestParameters["query"] != null) {
      queryParameters["query"] = requestParameters["query"];
    }

    const headerParameters: runtime.HTTPHeaders = {};

    const response = await this.request(
      {
        path: `/web/search`,
        method: "POST",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      WebSearchGet200ResponseFromJSON(jsonValue),
    );
  }

  /**
   * Search the content of a web page
   * Search web content
   */
  async webSearchPost(
    requestParameters: WebSearchPostRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<WebSearchGet200Response> {
    const response = await this.webSearchPostRaw(
      requestParameters,
      initOverrides,
    );
    return await response.value();
  }

  /**
   * Get the summary of a web page
   * Get web summary
   */
  async webSummaryGetRaw(
    requestParameters: WebSummaryGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<WebSummaryGet200Response>> {
    if (requestParameters["url"] == null) {
      throw new runtime.RequiredError(
        "url",
        'Required parameter "url" was null or undefined when calling webSummaryGet().',
      );
    }

    const queryParameters: any = {};

    if (requestParameters["url"] != null) {
      queryParameters["url"] = requestParameters["url"];
    }

    const headerParameters: runtime.HTTPHeaders = {};

    const response = await this.request(
      {
        path: `/web/summary`,
        method: "GET",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      WebSummaryGet200ResponseFromJSON(jsonValue),
    );
  }

  /**
   * Get the summary of a web page
   * Get web summary
   */
  async webSummaryGet(
    requestParameters: WebSummaryGetRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<WebSummaryGet200Response> {
    const response = await this.webSummaryGetRaw(
      requestParameters,
      initOverrides,
    );
    return await response.value();
  }

  /**
   * Get the summary of a web page
   * Get web summary
   */
  async webSummaryPostRaw(
    requestParameters: WebSummaryPostRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<runtime.ApiResponse<WebSummaryGet200Response>> {
    if (requestParameters["url"] == null) {
      throw new runtime.RequiredError(
        "url",
        'Required parameter "url" was null or undefined when calling webSummaryPost().',
      );
    }

    const queryParameters: any = {};

    if (requestParameters["url"] != null) {
      queryParameters["url"] = requestParameters["url"];
    }

    const headerParameters: runtime.HTTPHeaders = {};

    const response = await this.request(
      {
        path: `/web/summary`,
        method: "POST",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides,
    );

    return new runtime.JSONApiResponse(response, (jsonValue) =>
      WebSummaryGet200ResponseFromJSON(jsonValue),
    );
  }

  /**
   * Get the summary of a web page
   * Get web summary
   */
  async webSummaryPost(
    requestParameters: WebSummaryPostRequest,
    initOverrides?: RequestInit | runtime.InitOverrideFunction,
  ): Promise<WebSummaryGet200Response> {
    const response = await this.webSummaryPostRaw(
      requestParameters,
      initOverrides,
    );
    return await response.value();
  }
}
