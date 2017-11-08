import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { ClarityModule } from 'clarity-angular';
import { ToastyModule } from 'ng2-toasty';

import { SharedModule } from './shared/shared.module';
import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserAnimationsModule,
    HttpClientModule,
    BrowserModule,
    ClarityModule,
    FormsModule
  ],
  providers: [
    SharedModule
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
