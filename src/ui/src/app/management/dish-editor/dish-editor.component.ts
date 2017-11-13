import {Component, OnInit, Input, Output, EventEmitter} from '@angular/core';
import {FormGroup, FormControl, Validators} from '@angular/forms';
import {DISH_IMG_DEFAULT, IDish, DISH_TYPES} from '../../shared/interfaces/dish';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import { isNil } from 'lodash';
import {DishesService} from '../services/dishes.service';
import {WebHelperService} from '../../shared/services/web-helper.service';

/**
 * Dish item editor
 *
 * @example <app-dish-editor [dish]="someDish" (success)="onSuccess()"></app-dish-editor>
 */
@Component({
  selector: 'app-dish-editor',
  templateUrl: './dish-editor.component.html',
  styleUrls: ['./dish-editor.component.scss']
})
export class DishEditorComponent extends LoadStatusComponent implements OnInit {

  /**
   * Is control enabled
   * @type {boolean}
   */
  @Input() enabled = false;

  /**
   * Use compact mode (for modals)
   * @type {boolean}
   */
  @Input() compactMode = false;

  /**
   * Item to edit (only for existing item)
   * @type {any}
   */
  @Input() dish: IDish = null;

  @Output() success = new EventEmitter();

  @Output() dismiss = new EventEmitter();

  @Output() fail = new EventEmitter();

  @Output() loading = new EventEmitter();

  dishForm: FormGroup = null;

  imageLoading = false;

  private isNewDish = false;

  constructor(private dishes: DishesService, private helper: WebHelperService) {
    super();
  }

  get safePhotoUrl(): string {
    return !this.dish.photoUrl.length ? DISH_IMG_DEFAULT : this.dish.photoUrl;
  }

  get dishTypes() {
    return DISH_TYPES;
  }

  ngOnInit() {
    this.isNewDish = isNil(this.dish);
    if (this.isNewDish) {
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

  /**
   * Set form state (enabled/disabled)
   * @param {boolean} isEnabled
   */
  setFormState(isEnabled: boolean) {
    if (isEnabled) {
      this.dishForm.enable({onlySelf: false, emitEvent: false});
    } else {
      this.dishForm.disable({onlySelf: false, emitEvent: false});
    }
  }

  onImgUploadStart() {
    this.setFormState(false);
    this.imageLoading = true;
  }

  onImageUploaded(imgSrc: string) {
    this.dish.photoUrl = imgSrc;
    this.setFormState(true);
    this.imageLoading = false;
  }

  onImageUploadFail(err: string) {
    this.setFormState(true);
    this.imageLoading = false;
    window.alert(`Failed to upload the image: ${err || 'Unknown error'}`);
  }

  onSubmit(): boolean {
    if (this.dishForm.invalid || this.isLoading || this.imageLoading) {
      return false;
    }

    // Copy values from the form
    Object.assign(this.dish, this.dishForm.value);

    this.isLoading = true;
    this.loading.emit();

    if (this.isNewDish) {
      this.dishes.addDish(this.dish).subscribe(
        () => this.onSaveSuccess(),
        (err) => this.onSaveFail(err)
      );
    } else {
      this.dishes.updateDish(this.dish).subscribe(
        () => this.onSaveSuccess(),
        (err) => this.onSaveFail(err)
      );
    }

    return false;
  }

  private onSaveSuccess() {
    this.isLoaded = true;
    this.success.emit();
  }

  private onSaveFail(err: any) {
    this.isFailed = true;
    this.error = this.helper.extractResponseError(err);
    this.fail.emit();
  }

}
