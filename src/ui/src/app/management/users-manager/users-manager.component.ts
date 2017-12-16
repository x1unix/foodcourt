import { Component, OnInit } from '@angular/core';
import * as moment from 'moment';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {IUser} from '../../shared/interfaces/user';
import {SessionsService} from '../../shared/services/sessions.service';
import {UsersService} from '../../shared/services/users.service';
import {WebHelperService} from '../../shared/services/web-helper.service';
import {ResourceStatus} from '../../shared/helpers/resource-status';

const SUCCESS_ALERT_TIMEOUT = 5000;

const SERVED_DATE_FORMAT = 'YYYYMMDD';

@Component({
  selector: 'app-users-manager',
  templateUrl: './users-manager.component.html',
  styleUrls: ['./users-manager.component.scss']
})
export class UsersManagerComponent extends LoadStatusComponent implements OnInit {

  users: IUser[] = [];

  currentUserId: number = null;

  checkedIds: number[] = [];

  blockModal = false;

  showEditor = false;

  editableUser: IUser = null;

  editorOperation = 0;

  showSuccessMessage = false;

  deleteStatus = new ResourceStatus();

  todayDate: string = null;

  constructor(private usersService: UsersService, private session: SessionsService, private helper: WebHelperService) {
    super();
  }

  ngOnInit() {
    this.currentUserId = this.session.userId;
    this.todayDate = moment().format(SERVED_DATE_FORMAT);
    this.fetchUsers();
  }

  fetchUsers() {
    this.isLoading = true;

    // Just a tricky way to clear an array in JS (it's over 9000 ways to do it)
    this.checkedIds.length = 0;

    this.usersService.getAll().subscribe(
      (u) => this.onFetch(u),
      (e) => this.onFetchFail(e)
    );

  }

  onFetch(users: IUser[]) {
    this.users = users;
    this.isLoaded = true;
  }

  onFetchFail(err) {
    this.isFailed = true;
    this.error = this.helper.extractResponseError(err);
  }

  getGroupName(grpId: number): string {
    return this.usersService.getGroupName(grpId);
  }

  toggleUserCheck(uid: number) {
    const index = this.checkedIds.indexOf(uid);
    const isChecked = index > -1;

    if (isChecked) {
      this.checkedIds.splice(index, 1);
    } else {
      this.checkedIds.push(uid);
    }
  }

  isUserChecked(uid: number): boolean {
    return this.checkedIds.includes(uid);
  }

  /**
   * Editor progress event
   */
  onEditLoading() {
    this.blockModal = true;
  }

  /**
   * Editor work finish event
   * @param {boolean} isSuccess Is successful
   */
  onEditFinish(isSuccess: boolean) {
    this.blockModal = false;
    this.showEditor = false;
    this.editableUser = null;

    if (isSuccess === true) {
      this.showSuccessMessage = true;

      // Hide success message after timeout finish
      setTimeout(() => this.showSuccessMessage = false, SUCCESS_ALERT_TIMEOUT);

      // Refresh data
      this.fetchUsers();
    }
  }

  /**
   * Create new item in the editor
   */
  openItemCreator() {
    this.editableUser = null;
    this.showEditor = true;
    this.editorOperation = 0;
  }

  /**
   * Open item editor with the specified dish
   * @param {IUser} user Selected dish
   */
  editItem(user: IUser) {
    // Clear previous state
    this.editableUser = null;
    this.editorOperation = 1;
    this.editableUser = user;
    this.showEditor = true;
  }

  /**
   * Editor button close click event
   */
  onEditDismiss() {
    this.editableUser = null;
    this.blockModal = false;
    this.showEditor = false;
  }

  /**
   * Delete selected items
   */
  deleteSelected() {
    const confirm = window.confirm(`Delete ${this.checkedIds.length} item(s)?`);

    if (!confirm) {
      return;
    }

    this.deleteStatus.isLoading = true;

    // Hide delete operation result alert
    setTimeout(() => {
      this.deleteStatus.isIdle = true;
    }, 5000);

    this.usersService.removeMultiple(this.checkedIds).subscribe(
      () => {
        this.deleteStatus.isLoaded = true;
        this.fetchUsers();
      }, (e) => {
        this.deleteStatus.isFailed = true;
        this.deleteStatus.error = this.helper.extractResponseError(e);
      }
    );
  }

}
