import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-header-toolbar',
  templateUrl: './header-toolbar.component.html',
  styleUrls: ['./header-toolbar.component.scss']
})
export class HeaderToolbarComponent implements OnInit {
  @Input() label: string = null;

  @Input() smaller = false;

  constructor() { }

  ngOnInit() {
  }

}
