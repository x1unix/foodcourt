import { Component, OnInit } from '@angular/core';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {IUser} from '../../shared/interfaces/user';
import {SessionsService} from '../../shared/services/sessions.service';
import {UsersService} from '../../shared/services/users.service';
import {WebHelperService} from '../../shared/services/web-helper.service';

@Component({
  selector: 'app-users-manager',
  templateUrl: './users-manager.component.html',
  styleUrls: ['./users-manager.component.scss']
})
export class UsersManagerComponent extends LoadStatusComponent implements OnInit {

  users: IUser[] = [];

  currentUserId: number = null;

  checkedIds: number[] = [];

  constructor(private usersService: UsersService, private session: SessionsService, private helper: WebHelperService) {
    super();
  }

  ngOnInit() {
    this.currentUserId = this.session.userId;
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

}
