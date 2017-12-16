import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AuthComponent } from './auth.component';
import { SharedModule } from '../shared/shared.module';
import { ClarityModule } from 'clarity-angular';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ProfileEditorComponent } from './profile-editor/profile-editor.component';

@NgModule({
  imports: [
    CommonModule,
    RouterModule,
    SharedModule,
    FormsModule,
    ClarityModule,
    ReactiveFormsModule
  ],
  declarations: [AuthComponent, ProfileEditorComponent]
})
export class AuthModule { }
