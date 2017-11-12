import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { ClarityModule } from 'clarity-angular';
import { SharedModule } from '../shared/shared.module';

import { ManagementComponent } from './management.component';
import { ItemsCatalogComponent } from './items-catalog/items-catalog.component';
import { DishEditorComponent } from './dish-editor/dish-editor.component';
import {DishesService} from './services/dishes.service';
import { ImgPickerComponent } from './img-picker/img-picker.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    SharedModule,
    ClarityModule
  ],
  declarations: [ManagementComponent, ItemsCatalogComponent, DishEditorComponent, ImgPickerComponent],
  providers: [
    DishesService
  ]
})
export class ManagementModule { }
