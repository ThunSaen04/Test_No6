import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-root',
  imports: [FormsModule],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  private http = inject(HttpClient);
  protected product_code = ''
  private apiUrl = 'http://localhost:3000/api/barcodes';
  items = signal<any[]>([]);

  // items = signal<any[]>([
  //   {
  //     id: '1',
  //     product_code: '1234-4567-8909-8765',
  //     barcode_url: 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSPCMYq3TBikzP-dZCFAzfH7Fsu_49i4zI5fQ&s'
  //   },
  //   {
  //     id: '2',
  //     product_code: 'AK20-4567-8909-8765',
  //     barcode_url: 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSPCMYq3TBikzP-dZCFAzfH7Fsu_49i4zI5fQ&s'
  //   },
  //   {
  //     id: '3',
  //     product_code: 'SK2C-A29S-KGM2-8OII',
  //     barcode_url: 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSPCMYq3TBikzP-dZCFAzfH7Fsu_49i4zI5fQ&s'
  //   }
  // ]);

  ngOnInit() {
    this.fetchItems();
  }

  fetchItems() {

    this.http.get<any>(this.apiUrl).subscribe({
      next: (res) => {
        this.items.set(res.data ?? []);
      },
      error: (err) => {
        console.error('โหลดข้อมูลไม่สำเร็จ:', err);
      }
    });
  }

  addItem() {
    if (this.product_code.length !== 19) {
      alert('กรุณากรอกรหัสสินค้าให้ครบ 16 หลัก');
      return;
    }

    this.http.post<any>(this.apiUrl, { product_code: this.product_code }).subscribe({
      next: () => {
        this.fetchItems();
        this.product_code = '';
      },
      error: (err) => {
        console.error('เพิ่มสินค้าไม่สำเร็จ:', err);
        alert(err.error?.error || 'เพิ่มสินค้าไม่สำเร็จ');
      }
    });
  }

  //  deleteItem(id: string){
  //   const itemToDelete = this.items().find(item => item.id === id);

  //   if (!itemToDelete) return;

  //   const isConfirmed = window.confirm(`ต้องการลบข้อมูล รหัสสินค้า ${itemToDelete.product_code} หรือไม่?`);

  //   if (isConfirmed) {
  //     this.items.update((items) => items.filter((item) => item.id !== id));
  //   }
  //  }

  pendingDeleteId = signal<number | null>(null);

  getPendingCode() {
    const id = this.pendingDeleteId();
    return this.items().find(item => item.id === id)?.product_code || '';
  }

  deleteItem(id: number) {
    this.pendingDeleteId.set(id);
  }

  cancelDelete() {
    this.pendingDeleteId.set(null);
  }

  confirmDelete() {
    const id = this.pendingDeleteId();
    if (id) {
      this.http.delete(`${this.apiUrl}/${id}`).subscribe({
        next: () => {
          this.items.update(prev => prev.filter(item => item.id !== id));
          this.pendingDeleteId.set(null);
        },
        error: (err) => {
          console.error('ลบสินค้าไม่สำเร็จ:', err);
          alert(err.error?.error || 'ลบสินค้าไม่สำเร็จ');
        }
      });
    }
  }

  formatItemCode(event: any) {
    let value = event.target.value.replace(/[^a-zA-Z0-9]/g, '').toUpperCase();

    if (value.length > 16) {
      value = value.substring(0, 16);
    }

    let formattedValue = value.match(/.{1,4}/g)?.join('-') || '';

    this.product_code = formattedValue;
  }
}
