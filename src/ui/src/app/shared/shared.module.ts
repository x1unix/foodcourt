import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from 'clarity-angular';
import { ToastyModule } from 'ng2-toasty';

import {
  LocalStorageService,
  LoggerService,
  AuthService
} from './services';

import { TokenInterceptor } from './interceptors/token.interceptor';

/**
 * Module provides access to common app parts and services
 *
 * @export
 * @class SharedModule
 */
@NgModule({
  imports: [
    CommonModule,
    BrowserAnimationsModule,
    HttpClientModule,
    BrowserModule,
    ClarityModule,
    FormsModule
  ],
  declarations: [],
  providers: [
    LocalStorageService,
    LoggerService,
    AuthService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true
    }
  ],
  exports: [
    // ...FORM_DIRECTIVES,
    // ...COMPONENTS,
    // ...DIRECTIVES,
    // ...PIPES
  ]
})
export class SharedModule { }
