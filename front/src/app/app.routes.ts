import { Routes } from '@angular/router';
import { IndexComponent } from './components/index/index.component';
import { AttackComponent } from './components/attack/attack/attack.component';

export const routes: Routes = [
    { path: '', component: IndexComponent},
    { path: 'attack', component : AttackComponent}
];
