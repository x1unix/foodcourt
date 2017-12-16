import {Component, OnDestroy, OnInit} from '@angular/core';
import { isNil } from 'lodash';
import {LoadStatusComponent, ResourceStatus} from '../../shared/helpers';
import {IUser} from '../../shared/interfaces/user';
import {SessionsService, UsersService, WebHelperService} from '../../shared/services';
import {FormControl, FormGroup, Validators} from '@angular/forms';

const KEY_PASSWD = 'password';

@Component({
  selector: 'app-profile-editor',
  templateUrl: './profile-editor.component.html',
  styleUrls: ['./profile-editor.component.scss']
})
export class ProfileEditorComponent extends LoadStatusComponent implements OnInit, OnDestroy {

  /**
   * Data save status
   * @type {ResourceStatus}
   */
  saveStatus: ResourceStatus = null;

  /**
   * Current user id
   * @type {number}
   */
  userId: number = null;

  /**
   * Form
   * @type {null}
   */
  userForm: FormGroup = null;

  private userData: IUser = null;

  constructor(private users: UsersService, private session: SessionsService, private helper: WebHelperService) {
    super();
  }

  ngOnInit() {
    this.saveStatus = new ResourceStatus();
    this.userId = this.session.userId;
    this.fetchUserInfo();
  }

  ngOnDestroy() {
    this.saveStatus = null;
    this.userData = null;
    this.userForm = null;
  }

  fetchUserInfo() {
    this.users.findById(this.userId).subscribe(
      (user: IUser) => this.onFetch(user),
      (error) => this.onFetchError(this.helper.extractResponseError(error))
    );
  }

  onFetch(user: IUser) {
    delete user.password;

    this.userData = user;
    this.userForm = new FormGroup({
      firstName: new FormControl(this.userData.firstName, Validators.compose([
        Validators.required,
        Validators.maxLength(64)
      ])),
      lastName: new FormControl(this.userData.lastName, Validators.compose([
        Validators.required,
        Validators.maxLength(64)
      ])),
      email: new FormControl(this.userData.email, Validators.compose([
        Validators.email,
        Validators.maxLength(64),
        Validators.required
      ])),
      password: new FormControl('', Validators.maxLength(64))
    });
    this.isLoaded = true;
  }

  isInvalid(fieldName: string): boolean {
    const field = this.userForm.get(fieldName);
    return field.invalid && (field.dirty || field.touched);
  }

  onFetchError(error: string) {
    this.error = error;
    this.isFailed = true;
  }

  private getChanges(): IUser|any {
    const formValues = <IUser> this.userForm.value;
    const output: IUser|any = {};
    let keysChanged = 0;

    for (const key in formValues) {
      if (!formValues.hasOwnProperty(key)) {
        continue;
      }

      const val = formValues[key].trim();

      if (val === this.userData[key]) {
        continue;
      }

      if (key === KEY_PASSWD) {
        if (val.length === 0) {
          continue;
        }
      }

      output[key] = val;
      keysChanged++;
    }

    return keysChanged > 0 ? output : null;
  }

  onSubmit() {
    const changes = this.getChanges();

    if (isNil(changes)) {
      return;
    }

    this.saveStatus.isLoading = true;
    this.userForm.disable({emitEvent: false});

    this.users.updateUser(this.userId, changes).subscribe(
      () => {
        this.saveStatus.isLoaded = true;
        this.userForm.enable({emitEvent: false});
        setTimeout(() => this.saveStatus.isIdle = true, 3000);
      }, (e) => {
        this.saveStatus.error = this.helper.extractResponseError(e);
        this.userForm.enable({emitEvent: false});
        this.saveStatus.isFailed = true;
      }
    );
  }

}
