import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { IndexComponent } from "./components/index/index.component";

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, IndexComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'intruder-front';
}
