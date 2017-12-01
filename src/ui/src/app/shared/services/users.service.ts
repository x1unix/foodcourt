import { Injectable } from '@angular/core';

const GROUP_NAMES = [
  'Administrator',
  'Manager',
  'Customer'
];

@Injectable()
export class UsersService {

  constructor() { }

  /**
   * List of user groups
   * @returns {string[]}
   */
  get groups() {
    return GROUP_NAMES;
  }

  /**
   * Gets group name
   * @param {number} grpId Group ID
   * @returns {string}
   */
  getGroupName(grpId: number): string {
    return GROUP_NAMES[grpId];
  }

}
