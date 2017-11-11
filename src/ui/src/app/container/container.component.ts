import {Component, OnInit, ViewEncapsulation} from '@angular/core';
import {SessionsService} from '../shared/services/sessions.service';
import {AuthService} from '../shared/services/auth.service';
import {IUser} from '../shared/interfaces/user';
import {Router} from '@angular/router';

@Component({
  selector: 'app-container',
  templateUrl: './container.component.html',
  styleUrls: ['./container.component.scss']
})
export class ContainerComponent implements OnInit {

  currentUser: IUser = null;

  constructor(private session: SessionsService, private auth: AuthService, private router: Router) { }

  ngOnInit() {
    this.currentUser = this.session.currentUser;

    this.session.logout.subscribe(() => this.onLogout());
  }

  onLogoutClick() {
    this.auth.logout().subscribe((r) => {
      this.session.sessionLogout();
    });
  }

  onLogout() {
    this.currentUser = null;
    this.router.navigate(['/auth']);
  }

}
