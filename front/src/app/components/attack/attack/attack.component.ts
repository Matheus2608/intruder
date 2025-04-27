import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AttackMetadataComponent } from "../attack-metadata/attack-metadata.component";
import { RequestBodyComponent } from "../request-body/request-body.component";
import { AttackOutput } from '../../../interfaces/AttackOutput';
import { ResponseData } from '../../../interfaces/ResponseData';
import { ApiClientService } from '../../../services/api-client.service';

@Component({
  selector: 'app-attack',
  imports: [AttackMetadataComponent, RequestBodyComponent],
  templateUrl: './attack.component.html',
  styleUrl: './attack.component.css'
})
export class AttackComponent implements OnInit {
    elapsedTime: any;
    url: any;
    lenList: any;
    data: any | undefined;
    responses: ResponseData[] = [];

    constructor(private apiService : ApiClientService, private router: Router) {}

    ngOnInit() {
      this.apiService.getFakeApiResponse().subscribe({
        next: (response: AttackOutput) => {
        this.elapsedTime = response.totalElapsedTimeInMiliseconds;
        this.url = response.url;
        this.responses = response.responses;
      },
        error: (error) => {
          console.error('Error fetching data:', error);
      }});
    }

//     toggleDetails(row) {
//         if (row.classList.contains('details-row')) return;
//         const detailsRow = row.nextElementSibling;
//         if (detailsRow && detailsRow.classList.contains('details-row')) {
//             detailsRow.style.display = detailsRow.style.display === "none" ? "table-row" : "none";
//         }
//     }

//     activateTab(button) {
//       const container = button.closest('.request-response-container');
//       if (!container) return;
//       const targetClass = button.getAttribute('data-target');
//       container.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
//       container.querySelectorAll('.tab-content').forEach(content => content.style.display = 'none');
//       button.classList.add('active');
//       const targetContent = container.querySelector(`.${targetClass}`);
//       if (targetContent) {
//           targetContent.style.display = 'block';
//       }
//   }

//     setupTabListeners() {
//       document.querySelectorAll('.tab-button').forEach(button => {
//           button.removeEventListener('click', handleTabClick);
//           button.addEventListener('click', handleTabClick);
//       });
//   }

//     handleTabClick(event) {
//       activateTab(event.currentTarget);
//   }


//   // --- Sorting Logic ---
//   let currentSortKey = 'RequestId';
//   let currentSortDir = 'asc';

//   function getCellValue(row, sortKey) {
//       const cell = row.querySelector(`td[data-type="${sortKey}"]`);
//       if (!cell) return '';
//       let value = cell.textContent.trim();
//       switch (sortKey) {
//           case 'RequestId':
//           case 'StatusCode':
//           case 'TimeElapsed':
//           case 'Length':
//               const num = parseInt(value, 10);
//               // *** POTENCIAL PROBLEMA AQUI: Se status é N/A, parseInt dá NaN. Comparar NaN falha. ***
//               // *** Vamos tratar NaN explicitamente. ***
//               return isNaN(num) ? (sortKey === 'StatusCode' ? -1 : 0) : num; // Treat N/A status as -1, others as 0? or use a very small/large number? -1 for Status is okay. 0 for others might mess order.
//               // Melhor: retornar um valor consistente para NaN, como Number.MIN_SAFE_INTEGER para asc ou Number.MAX_SAFE_INTEGER para desc, dependendo da direção.
//               // return isNaN(num) ? (currentSortDir === 'asc' ? Number.MIN_SAFE_INTEGER : Number.MAX_SAFE_INTEGER) : num; // Mais robusto?
//           case 'Err':
//               return value.toLowerCase() === 'yes' ? 1 : 0;
//           case 'Payload':
//           default:
//               return value.toLowerCase();
//       }
//   }


//   function sortTable(sortKey, headerCell) {
//       const table = document.getElementById('resultsTable');
//       const tbody = table.querySelector('tbody');
//       if (!tbody) return;

//       // *** POTENCIAL PROBLEMA AQUI: Pegar TODAS as TRs, incluindo as de detalhes, pode quebrar a lógica ***
//       // const rows = Array.from(tbody.querySelectorAll('tr:not(.details-row)')); // Get only main data rows

//       // *** CORREÇÃO: Vamos pegar os pares Main+Details ANTES de ordenar ***
//       const allRows = Array.from(tbody.querySelectorAll('tr'));
//       const rowPairs = [];
//       for (let i = 0; i < allRows.length; i += 2) {
//           if (allRows[i] && !allRows[i].classList.contains('details-row') && // É uma linha principal
//               allRows[i+1] && allRows[i+1].classList.contains('details-row')) { // E a próxima é de detalhes
//               rowPairs.push({ main: allRows[i], details: allRows[i+1] });
//           } else if (allRows[i]) {
//               // Linha principal sem linha de detalhes? Adicionar apenas ela.
//               // Ou pode indicar um problema no HTML gerado pelo Go template se sempre deveriam ter pares.
//               console.warn("Found main row without subsequent details row at index", i, allRows[i]);
//               // Vamos incluir apenas a linha principal por enquanto se não tiver par
//               // rowPairs.push({ main: allRows[i], details: null });
//           }
//       }
//       // *** FIM CORREÇÃO ***


//       // Determine sort direction
//       if (sortKey === currentSortKey) {
//           currentSortDir = currentSortDir === 'asc' ? 'desc' : 'asc';
//       } else {
//           currentSortDir = 'asc';
//           table.querySelectorAll('thead th').forEach(th => th.classList.remove('sort-asc', 'sort-desc'));
//       }
//       currentSortKey = sortKey;

//       // Update header indicator
//       if (headerCell) {
//           headerCell.classList.remove('sort-asc', 'sort-desc');
//           headerCell.classList.add(currentSortDir === 'asc' ? 'sort-asc' : 'sort-desc');
//       }

//       // Sort the ROW PAIRS based on the main row's value
//       rowPairs.sort((a, b) => {
//           const aValue = getCellValue(a.main, sortKey); // Usa a.main
//           const bValue = getCellValue(b.main, sortKey); // Usa b.main

//           let comparison = 0;
//           // *** POTENCIAL PROBLEMA: Comparação direta pode falhar para NaN ou tipos mistos ***
//           // Vamos refinar a comparação
//           if (typeof aValue === 'number' && typeof bValue === 'number') {
//               comparison = aValue - bValue; // Comparação numérica direta
//           } else {
//               // Comparação como string se não forem ambos números
//               const stringA = String(aValue);
//               const stringB = String(bValue);
//               if (stringA < stringB) { comparison = -1; }
//               else if (stringA > stringB) { comparison = 1; }
//           }


//           // if (aValue < bValue) { // Comparação antiga
//           //     comparison = -1;
//           // } else if (aValue > bValue) {
//           //     comparison = 1;
//           // }

//           return currentSortDir === 'asc' ? comparison : comparison * -1;
//       });

//       // Re-append row pairs in sorted order
//       tbody.innerHTML = ''; // Clear existing content
//       rowPairs.forEach(pair => {
//           tbody.appendChild(pair.main);
//           if (pair.details) { // Anexa a linha de detalhes SE ela existir
//             tbody.appendChild(pair.details);
//           }
//       });

//       // Re-attach event listeners for tabs after sorting
//       setupTabListeners();
//   }


//   // --- Multi-level Sort Controls (Optional - Keep or Remove) ---
//   function addOrder() { /* ... keep original addOrder JS ... */ }
//   function removeOrder(button) { /* ... keep original removeOrder JS ... */ }
//   function applySorting() { /* Implement if using multi-level sort */ }


//   // --- Initial Setup ---
//   document.addEventListener('DOMContentLoaded', () => {
//       setupTabListeners(); // Setup listeners for initial tabs
//       // Apply default sort on load - ESTA É A LINHA CRÍTICA
//       const defaultSortHeader = document.querySelector(`thead th[data-sort-key="${currentSortKey}"]`);
//       if (defaultSortHeader) { // Verifica se o header existe
//           try {
//               sortTable(currentSortKey, defaultSortHeader); // Tenta ordenar
//           } catch (error) {
//               console.error("Error during initial table sort:", error);
//               // Se falhar, a tabela pode ficar vazia ou parcialmente renderizada.
//               // Poderíamos tentar exibir a tabela sem ordenar como fallback?
//           }
//       } else {
//           console.warn("Default sort header not found for key:", currentSortKey);
//       }

//   });
}
