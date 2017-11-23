import { Injectable } from '@angular/core';
import {HttpErrorResponse} from '@angular/common/http';
import { isObject, isString } from 'lodash';
import {IMessage} from '../interfaces/message';

@Injectable()
export class WebHelperService {

  constructor() { }

  /**
   * extract response error
   * @param {HttpErrorResponse} err
   * @returns {string}
   */
  extractResponseError(err: HttpErrorResponse): string {
    if (isString(err.error)) {
      return err.error;
    }

    if (isObject(err.error)) {
      const msg: IMessage = <IMessage> err.error;
      return `${err.status} ${msg.msg}`;
    }

    return `${err.status} ${err.statusText}`;
  }

}
