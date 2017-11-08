import {Component, OnInit, ViewEncapsulation} from '@angular/core';

@Component({
  selector: 'app-food-card',
  templateUrl: './food-card.component.html',
  styleUrls: ['./food-card.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class FoodCardComponent implements OnInit {

  constructor() { }

  ngOnInit() {
  }

}
