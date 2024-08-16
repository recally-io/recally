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
 * @interface JinasearcherContent
 */
export interface JinasearcherContent {
    /**
     * 
     * @type {string}
     * @memberof JinasearcherContent
     */
    content?: string;
    /**
     * 
     * @type {string}
     * @memberof JinasearcherContent
     */
    description?: string;
    /**
     * 
     * @type {string}
     * @memberof JinasearcherContent
     */
    title?: string;
    /**
     * 
     * @type {string}
     * @memberof JinasearcherContent
     */
    url?: string;
}

/**
 * Check if a given object implements the JinasearcherContent interface.
 */
export function instanceOfJinasearcherContent(value: object): value is JinasearcherContent {
    return true;
}

export function JinasearcherContentFromJSON(json: any): JinasearcherContent {
    return JinasearcherContentFromJSONTyped(json, false);
}

export function JinasearcherContentFromJSONTyped(json: any, ignoreDiscriminator: boolean): JinasearcherContent {
    if (json == null) {
        return json;
    }
    return {
        
        'content': json['content'] == null ? undefined : json['content'],
        'description': json['description'] == null ? undefined : json['description'],
        'title': json['title'] == null ? undefined : json['title'],
        'url': json['url'] == null ? undefined : json['url'],
    };
}

export function JinasearcherContentToJSON(value?: JinasearcherContent | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'content': value['content'],
        'description': value['description'],
        'title': value['title'],
        'url': value['url'],
    };
}

