import {Component, OnDestroy, OnInit, ViewEncapsulation} from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {HttpErrorResponse} from '@angular/common/http';

import { isNil } from 'lodash';

import {AuthService, SessionsService, WebHelperService} from '../shared/services';
import {LoadStatusComponent} from '../shared/helpers';
import {IAuthCredentials} from '../shared/interfaces/auth-credentials';
import {IAuthSession} from '../shared/interfaces/auth-session';
import {Subscription} from 'rxjs/Subscription';

@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  styleUrls: ['./auth.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class AuthComponent extends LoadStatusComponent implements OnInit, OnDestroy {

  title = 'Limelight FoodCourt';

  sessionExpired = false;

  data: IAuthCredentials = null;

  routeParams$: Subscription = null;

  constructor(
    private session: SessionsService,
    private auth: AuthService,
    private hlp: WebHelperService,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    super();
  }

  ngOnInit() {
    // Redirect if logged in
    if (this.session.isAuthorized) {
      this.router.navigate(['/']);
    }

    this.data = {
      email: '',
      password: ''
    };

    // Check route params
    this.routeParams$ = this.route.queryParams.subscribe(params => {
      this.sessionExpired = !isNil(params['error']);
      this.error = 'Session expired, please log in';
    });
  }

  tryLogin(event: Event) {
    event.preventDefault();
    this.isLoading = true;
    this.auth.login(this.data).subscribe(
      (data: IAuthSession) => {
        this.isLoaded = true;
        this.session.loadSession(data);
        this.router.navigate(['/']);
      },
      (err: HttpErrorResponse) => {
        this.isFailed = true;
        this.error = this.hlp.extractResponseError(err);
      }
    );
  }

  ngOnDestroy() {
    this.routeParams$.unsubscribe();
  }

}
