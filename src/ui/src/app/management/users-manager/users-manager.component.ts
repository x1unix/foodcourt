import { Component, OnInit } from '@angular/core';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {IUser} from '../../shared/interfaces/user';
import {SessionsService} from '../../shared/services/sessions.service';
import {UsersService} from '../../shared/services/users.service';

@Component({
  selector: 'app-users-manager',
  templateUrl: './users-manager.component.html',
  styleUrls: ['./users-manager.component.scss']
})
export class UsersManagerComponent extends LoadStatusComponent implements OnInit {

  users: IUser[] = [
    {
      id: 1,
      email: 'admin@llnw.com',
      firstName: 'Super',
      lastName: 'Admin',
      level: 0
    },
    {
      id: 2,
      email: 'dsedchenko@llnw.com',
      firstName: 'Denis',
      lastName: 'Sedchenko',
      level: 2
    }
  ];

  currentUserId: number = null;

  checkedIds: number[] = [];

  constructor(private usersService: UsersService, private session: SessionsService) {
    super();
  }

  ngOnInit() {
    this.currentUserId = this.session.userId;
  }

  fetchUsers() {

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
