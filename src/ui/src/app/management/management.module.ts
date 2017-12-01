import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { ClarityModule } from 'clarity-angular';
import { SharedModule } from '../shared/shared.module';

import {DndModule} from 'ng2-dnd';

import { ManagementComponent } from './management.component';
import { ItemsCatalogComponent } from './items-catalog/items-catalog.component';
import { DishEditorComponent } from './dish-editor/dish-editor.component';
import {DishesService} from './services/dishes.service';
import { ImgPickerComponent } from './img-picker/img-picker.component';
import { MenuEditorComponent } from './menu-editor/menu-editor.component';
import { DishesListComponent } from './dishes-list/dishes-list.component';
import { UsersManagerComponent } from './users-manager/users-manager.component';
import { UserEditorComponent } from './user-editor/user-editor.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    SharedModule,
    ClarityModule,
    DndModule.forRoot()
  ],
  declarations: [
    ManagementComponent,
    ItemsCatalogComponent,
    DishEditorComponent,
    ImgPickerComponent,
    MenuEditorComponent,
    DishesListComponent,
    UsersManagerComponent,
    UserEditorComponent
  ],
  providers: [
    DishesService
  ]
})
export class ManagementModule { }
