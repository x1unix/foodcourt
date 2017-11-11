import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {HttpClientModule} from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { ToastyModule } from 'ng2-toasty';

import { ManagementModule} from './management/management.module';
import { SharedModule } from './shared/shared.module';
import { ContainerModule } from './container/container.module';
import { AuthModule } from './auth/auth.module';
import { AppRoutingModule } from './app.routing';

import { AppComponent } from './app.component';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserAnimationsModule,
    HttpClientModule,
    BrowserModule,
    RouterModule,
    FormsModule,
    ToastyModule,
    AuthModule,
    ContainerModule,
    SharedModule,
    AppRoutingModule,
    ManagementModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
