import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse
} from '@angular/common/http';
import {Observable} from 'rxjs/Rx';
import { SessionsService, LoggerService, URL_LOGIN } from '../services';

const WHITELIST = [
  URL_LOGIN
];

/**
 * Interceptor handles token and response routine for API communication
 */
@Injectable()
export class TokenInterceptor implements HttpInterceptor {
  constructor(private auth: SessionsService, private log: LoggerService) {}
  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {

    // Check if necessary to add token
    const addToken = (!WHITELIST.includes(request.url)) || this.auth.isAuthorized;

    if (addToken) {
      // Add token to request's query params
      request = request.clone({
        setParams: {
          'token': this.auth.token
        }
      });
    }

    // Handle HTTP errors
    return next.handle(request)
      .map((event: HttpEvent<any>) => event)
      .catch((err: any, caught) => {
        if (err instanceof HttpErrorResponse) {

          // Emit event if session is expired
          if ((err.status === 403) || (err.status === 401) && !WHITELIST.includes(request.url)) {
            this.auth.purgeSession(true);
          }

          return Observable.throw(err);
        }
      });
  }
}
