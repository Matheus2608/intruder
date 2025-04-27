import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AttackMetadataComponent } from './attack-metadata.component';

describe('AttackMetadataComponent', () => {
  let component: AttackMetadataComponent;
  let fixture: ComponentFixture<AttackMetadataComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AttackMetadataComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AttackMetadataComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
