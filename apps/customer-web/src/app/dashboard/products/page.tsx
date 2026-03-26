import { Card, CardContent, CardHeader, CardTitle } from '@repo/ui/card';
import { Button } from '@repo/ui/button';
import { Input } from '@repo/ui/input';
import { Search, Plus, Filter } from 'lucide-react';

// Mock data for products
const products = [
  {
    id: 1,
    name: 'Premium T-Shirt',
    category: 'Clothing',
    price: 29.99,
    stock: 45,
    image: 'https://via.placeholder.com/80',
  },
  {
    id: 2,
    name: 'Wireless Earbuds',
    category: 'Electronics',
    price: 89.99,
    stock: 23,
    image: 'https://via.placeholder.com/80',
  },
  {
    id: 3,
    name: 'Leather Wallet',
    category: 'Accessories',
    price: 49.99,
    stock: 12,
    image: 'https://via.placeholder.com/80',
  },
  {
    id: 4,
    name: 'Smart Watch',
    category: 'Electronics',
    price: 199.99,
    stock: 8,
    image: 'https://via.placeholder.com/80',
  },
  {
    id: 5,
    name: 'Running Shoes',
    category: 'Footwear',
    price: 79.99,
    stock: 15,
    image: 'https://via.placeholder.com/80',
  },
  {
    id: 6,
    name: 'Backpack',
    category: 'Accessories',
    price: 59.99,
    stock: 20,
    image: 'https://via.placeholder.com/80',
  },
];

export default function ProductsPage() {
  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between space-y-4 sm:flex-row sm:items-center">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Products</h1>
          <p className="text-muted-foreground">Manage your products and inventory</p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" />
          Add Product
        </Button>
      </div>

      <Card>
        <CardHeader className="flex flex-col space-y-4 sm:flex-row sm:items-center sm:justify-between sm:space-y-0">
          <div className="relative w-full max-w-sm">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="search"
              placeholder="Search products..."
              className="pl-8"
            />
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" className="h-9">
              <Filter className="mr-2 h-4 w-4" />
              Filter
            </Button>
            <Button variant="outline" size="sm" className="h-9">
              Export
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {products.map((product) => (
              <div key={product.id} className="group relative overflow-hidden rounded-lg border">
                <div className="aspect-square bg-gray-100 p-4">
                  <img
                    src={product.image}
                    alt={product.name}
                    className="h-full w-full object-contain"
                  />
                </div>
                <div className="p-4">
                  <h3 className="font-medium">{product.name}</h3>
                  <p className="text-sm text-muted-foreground">{product.category}</p>
                  <div className="mt-2 flex items-center justify-between">
                    <span className="font-bold">${product.price.toFixed(2)}</span>
                    <span className={`text-sm ${
                      product.stock > 10 ? 'text-green-600' : 'text-amber-600'
                    }`}>
                      {product.stock} in stock
                    </span>
                  </div>
                  <div className="mt-4 flex space-x-2">
                    <Button variant="outline" size="sm" className="flex-1">
                      Edit
                    </Button>
                    <Button variant="outline" size="sm">
                      View
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
