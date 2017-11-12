import { Component, OnInit, Input } from '@angular/core';
import {FormGroup, FormControl, Validators} from '@angular/forms';
import {DISH_IMG_DEFAULT, IDish, DISH_TYPES} from '../../shared/interfaces/dish';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import { isNil } from 'lodash';

@Component({
  selector: 'app-dish-editor',
  templateUrl: './dish-editor.component.html',
  styleUrls: ['./dish-editor.component.scss']
})
export class DishEditorComponent extends LoadStatusComponent implements OnInit {

  @Input() enabled = false;

  @Input() dish: IDish = null;

  dishForm: FormGroup = null;

  constructor() {
    super();
  }

  get safePhotoUrl(): string {
    return !this.dish.photoUrl.length ? DISH_IMG_DEFAULT : this.dish.photoUrl;
  }

  get dishTypes() {
    return DISH_TYPES;
  }

  ngOnInit() {
    if (isNil(this.dish)) {
      this.dish = {
        label: '',
        description: '',
        type: 0,
        photoUrl: ''
      };
    }

    this.dishForm = new FormGroup({
      label: new FormControl(this.dish.label, Validators.required),
      description: new FormControl(this.dish.description),
      type: new FormControl(this.dish.type)
    });
  }

  isInvalid(fieldName: string): boolean {
    const field = this.dishForm.get(fieldName);
    return field.invalid && (field.dirty || field.touched);
  }

  onSubmit() {

  }

}
