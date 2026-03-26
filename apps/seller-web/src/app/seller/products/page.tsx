'use client';
import React, { useState } from 'react';
import Card from '@/components/ui/Card';
import Table from '@/components/ui/Table';
import Pagination from '@/components/ui/Pagination';
import { Button } from '@repo/ui/button';
import { PlusCircle } from 'lucide-react';
import ActionDropdown from '@/components/ui/ActionDropdown';
import { Edit, Trash2 } from 'lucide-react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
  category: string;
}

const mockProducts: Product[] = [
  { id: '1', name: 'Laptop', price: 1200, stock: 50, category: 'Electronics' },
  { id: '2', name: 'Mouse', price: 25, stock: 200, category: 'Electronics' },
  { id: '3', name: 'Keyboard', price: 75, stock: 100, category: 'Electronics' },
  { id: '4', name: 'Monitor', price: 300, stock: 30, category: 'Electronics' },
  { id: '5', name: 'Webcam', price: 50, stock: 150, category: 'Electronics' },
  { id: '6', name: 'Desk Chair', price: 150, stock: 40, category: 'Furniture' },
  { id: '7', name: 'Bookshelf', price: 80, stock: 60, category: 'Furniture' },
  { id: '8', name: 'Coffee Maker', price: 90, stock: 70, category: 'Appliances' },
  { id: '9', name: 'Blender', price: 60, stock: 90, category: 'Appliances' },
  { id: '10', name: 'Toaster', price: 40, stock: 120, category: 'Appliances' },
];

const ProductsPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const productsPerPage = 5;
  const router = useRouter();

  const indexOfLastProduct = currentPage * productsPerPage;
  const indexOfFirstProduct = indexOfLastProduct - productsPerPage;
  const currentProducts = mockProducts.slice(indexOfFirstProduct, indexOfLastProduct);

  const totalPages = Math.ceil(mockProducts.length / productsPerPage);

  const handlePageChange = (pageNumber: number) => {
    setCurrentPage(pageNumber);
  };

  const handleEdit = (product: Product) => {
    router.push(`/seller/products/edit/${product.id}`);
  };

  const handleDelete = (product: Product) => {
    alert(`Delete product: ${product.name}`);
    // TODO: Implement actual delete logic
  };

  const columns = [
    { header: 'Product Name', accessor: 'name' },
    { header: 'Price', accessor: (row: Product) => `${row.price.toFixed(2)}` },
    { header: 'Stock', accessor: 'stock' },
    { header: 'Category', accessor: 'category' },
    {
      header: 'Actions',
      accessor: (row: Product) => (
        <ActionDropdown
          actions={[
            { label: 'Edit', icon: Edit, onClick: () => handleEdit(row) },
            { label: 'Delete', icon: Trash2, onClick: () => handleDelete(row) },
          ]}
        />
      ),
    },
  ];

  return (
    <div className="p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Products</h2>
        <Link href="/seller/products/add">
          <Button leftIcon={<PlusCircle className="h-5 w-5" />}>
            Add New Product
          </Button>
        </Link>
      </div>

      <Card className="mb-4">
        <Table columns={columns} data={currentProducts} />
      </Card>

      <div className="flex justify-center">
        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={handlePageChange}
        />
      </div>
    </div>
  );
};

export default ProductsPage;
