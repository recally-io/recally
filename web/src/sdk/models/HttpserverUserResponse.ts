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

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface HttpserverUserResponse
 */
export interface HttpserverUserResponse {
    /**
     * 
     * @type {string}
     * @memberof HttpserverUserResponse
     */
    email?: string;
    /**
     * 
     * @type {string}
     * @memberof HttpserverUserResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof HttpserverUserResponse
     */
    username?: string;
}

/**
 * Check if a given object implements the HttpserverUserResponse interface.
 */
export function instanceOfHttpserverUserResponse(value: object): value is HttpserverUserResponse {
    return true;
}

export function HttpserverUserResponseFromJSON(json: any): HttpserverUserResponse {
    return HttpserverUserResponseFromJSONTyped(json, false);
}

export function HttpserverUserResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): HttpserverUserResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'email': json['email'] == null ? undefined : json['email'],
        'id': json['id'] == null ? undefined : json['id'],
        'username': json['username'] == null ? undefined : json['username'],
    };
}

export function HttpserverUserResponseToJSON(value?: HttpserverUserResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'email': value['email'],
        'id': value['id'],
        'username': value['username'],
    };
}

