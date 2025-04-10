import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PayloadAreaComponent } from './payload-area.component';

describe('PayloadAreaComponent', () => {
  let component: PayloadAreaComponent;
  let fixture: ComponentFixture<PayloadAreaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PayloadAreaComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PayloadAreaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
