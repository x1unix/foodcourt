import { Injectable } from '@angular/core';
import { Router, CanActivate } from '@angular/router';
import { SessionsService } from '../services/sessions.service';

/**
 * Route guard that checks if user is authorized
 *
 * @export
 * @class AdminGuard
 * @implements {CanActivate}
 */
@Injectable()
export class AdminGuard implements CanActivate {
  constructor(private auth: SessionsService, private router: Router) {}

  /**
   * Page access checker
   *
   * @returns
   * @memberof LoggedInGuard
   */
  canActivate() {
    const isAdmin = this.auth.currentUser.level === 0;
    if (!isAdmin) {
      this.router.navigate(['/']);
    }

    return isAdmin;
  }
}
