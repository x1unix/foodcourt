import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import {CUSTOM_ELEMENTS_SCHEMA, NgModule, NO_ERRORS_SCHEMA} from '@angular/core';
import { CommonModule } from '@angular/common';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from 'clarity-angular';
import { ToastyModule } from 'ng2-toasty';
import { RouterModule } from '@angular/router';

import {
  LocalStorageService,
  LoggerService,
  SessionsService,
  AuthService,
  WebHelperService
} from './services';

import { TokenInterceptor } from './interceptors/token.interceptor';
import { LoggedInGuard } from './guards/logged-in.guard';
import { FoodCardComponent } from './components/food-card/food-card.component';
import { RatingComponent } from './components/rating/rating.component';
/**
 * Module provides access to common app parts and services
 *
 * @export
 * @class SharedModule
 */
@NgModule({
  imports: [
    ToastyModule,
    CommonModule,
    BrowserAnimationsModule,
    HttpClientModule,
    BrowserModule,
    ClarityModule,
    RouterModule,
    FormsModule
  ],
  declarations: [FoodCardComponent, RatingComponent],
  providers: [
    LocalStorageService,
    LoggerService,
    SessionsService,
    LoggedInGuard,
    AuthService,
    WebHelperService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true
    }
  ],
  schemas: [
    CUSTOM_ELEMENTS_SCHEMA
  ],
  exports: [
    FoodCardComponent
    // ...FORM_DIRECTIVES,
    // ...COMPONENTS,
    // ...DIRECTIVES,
    // ...PIPES
  ]
})
export class SharedModule { }
