import {Component, OnInit} from '@angular/core';
import { Title } from '@angular/platform-browser';
import {ActivatedRoute, NavigationEnd, Router} from '@angular/router';

import { isNil } from 'lodash';

const APP_TITLE = 'Food Court';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'app';

  constructor(
    private titleService: Title,
    private activatedRoute: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit() {
    // Subscribe for router navigation changes
    this.router.events
      .filter(event => event instanceof NavigationEnd)
      .map(() => this.activatedRoute)
      .map(route => {
        while (route.firstChild) { route = route.firstChild; }
        return route;
      })
      .filter(route => route.outlet === 'primary')
      .mergeMap(route => route.data)
      .subscribe((event) => {
        // Apply page title for the new route
        const title = isNil(event['title']) ? APP_TITLE : `${APP_TITLE} | ${event['title']}`;
        this.titleService.setTitle(title);
      });
  }
}
