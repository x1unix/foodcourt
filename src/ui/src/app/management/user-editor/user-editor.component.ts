import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { isNil, isEmpty } from 'lodash';
import {IUser} from '../../shared/interfaces/user';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {UsersService, Level} from '../../shared/services/users.service';
import {SessionsService} from '../../shared/services/sessions.service';
import {WebHelperService} from '../../shared/services/web-helper.service';


@Component({
  selector: 'app-user-editor',
  templateUrl: './user-editor.component.html',
  styleUrls: ['./user-editor.component.scss']
})
export class UserEditorComponent extends LoadStatusComponent implements OnInit {

  @Input() enabled = false;

  @Input() user: IUser = null;

  @Output() success = new EventEmitter();

  @Output() fail = new EventEmitter();

  @Output() loading = new EventEmitter();

  userForm: FormGroup = null;

  userId: number = null;

  currentUserId: number = null;

  userExists = false;

  constructor(private users: UsersService, private session: SessionsService, private helper: WebHelperService) {
    super();
  }

  get groups() {
    return this.users.groups;
  }

  isInvalid(fieldName: string): boolean {
    const field = this.userForm.get(fieldName);
    return field.invalid && (field.dirty || field.touched);
  }

  ngOnInit() {
    this.currentUserId = this.session.userId;

    this.userExists = !isNil(this.user);

    if (this.userExists) {
      this.userId = this.user.id;
    }

    if (this.userExists) {
      this.userForm = new FormGroup({
        firstName: new FormControl(this.user.firstName),
        lastName: new FormControl(this.user.lastName),
        email: new FormControl(this.user.email, Validators.email),
        password: new FormControl('')
        // level: new FormControl(this.user.level)
      });

      // if (this.userId === this.currentUserId) {
      //   this.userForm.get('level').disable({onlySelf: false, emitEvent: false});
      // }

    } else {
      this.userForm = new FormGroup({
        firstName: new FormControl('', Validators.required),
        lastName: new FormControl('', Validators.required),
        email: new FormControl('', Validators.compose([Validators.required, Validators.email])),
        password: new FormControl('', Validators.required),
        level: new FormControl(Level.Customer, Validators.required)
      });
    }
  }

  /**
   * Set form state (enabled/disabled)
   * @param {boolean} isEnabled
   */
  setFormState(isEnabled: boolean) {
    if (isEnabled) {
      this.userForm.enable({onlySelf: false, emitEvent: false});
    } else {
      this.userForm.disable({onlySelf: false, emitEvent: false});
    }
  }

  /**
   * Get changed values in form
   * @returns {IUser | any}
   */
  private getFormChanges() {
    const output: IUser|any = {};

    const data: IUser = this.userForm.value;
    const keys = Object.keys(data);

    keys.forEach((key) => {
      switch (key) {
        case 'id':
          return;

        case 'password':
          if (data.password.trim().length > 0) {
            output.password = data.password;
          }
          return;

        case 'level':
          if (this.user.level !== data.level) {
            output.level = Number(data.level);
          }
          return;

        default:
          const was = this.user[key].trim();
          const became = data[key].trim();

          if (was !== became) {
            output[key] = became;
          }
          return;
      }
    });

    return output;
  }

  onSubmit() {
    const data = this.userExists ? this.getFormChanges() : this.userForm.value;
    const req = this.userExists ? this.users.updateUser(this.userId, data) : this.users.addUser(data);


    // Workaround to pass only numbers as user level
    data.level = Number(data.level);

    if (isNaN(data.level)) {
      delete data.level;
    }

    if (isEmpty(data)) {
      this.isFailed = true;
      this.error = 'The form is empty';
      return;
    }

    this.isLoading = true;
    this.setFormState(false);
    this.loading.emit();

    req.subscribe(
      () => this.onSuccess(),
      (err) => this.onFail(err)
    );
  }

  onSuccess() {
    this.isLoaded = true;
    this.setFormState(true);
    this.success.emit();
  }

  onFail(err) {
    this.isFailed = true;
    this.setFormState(true);
    this.error = this.helper.extractResponseError(err);
    this.fail.emit();
  }

}
