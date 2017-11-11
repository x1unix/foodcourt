import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RetryAlertComponent } from './retry-alert.component';

describe('RetryAlertComponent', () => {
  let component: RetryAlertComponent;
  let fixture: ComponentFixture<RetryAlertComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RetryAlertComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RetryAlertComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
