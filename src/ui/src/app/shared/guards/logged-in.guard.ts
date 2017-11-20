import { Injectable } from '@angular/core';
import { Router, CanActivate } from '@angular/router';
import { SessionsService } from '../services/sessions.service';

/**
 * Route guard that checks if user is authorized
 *
 * @export
 * @class LoggedInGuard
 * @implements {CanActivate}
 */
@Injectable()
export class LoggedInGuard implements CanActivate {
  constructor(private auth: SessionsService, private router: Router) {}

  /**
   * Page access checker
   *
   * @returns
   * @memberof LoggedInGuard
   */
  canActivate() {
    const isAuthorized = this.auth.isAuthorized;
    if (!isAuthorized) {
      this.router.navigate(['/auth']);
    }

    return isAuthorized;
  }
}
