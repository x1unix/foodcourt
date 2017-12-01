import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { isNil } from 'lodash';
import {IUser} from '../../shared/interfaces/user';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {UsersService} from '../../shared/services/users.service';
import {SessionsService} from '../../shared/services/sessions.service';

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

  @Output() dismiss = new EventEmitter();

  userForm: FormGroup = null;

  userId: number = null;

  currentUserId: number = null;

  userExists = false;

  constructor(private users: UsersService, private session: SessionsService) {
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


    this.userForm = new FormGroup({
      firstName: new FormControl(this.userExists ? this.user.firstName : '', Validators.required),
      lastName: new FormControl(this.userExists ? this.user.lastName : '', Validators.required),
      // type: new FormControl(this.dish.type)
    });
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

  onSubmit() {

  }

}
