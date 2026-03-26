"use client";

import React, { useEffect, useState } from 'react';
import { Form, Field } from 'react-final-form';
import { Button } from '@repo/ui/button';
import Card from '@/components/ui/Card';
import { useRouter } from 'next/navigation';

interface EditProductFormData {
  name: string;
  description: string;
  price: number;
  stock: number;
  category: string;
}

const mockProducts = [
  { id: '1', name: 'Laptop', price: 1200, stock: 50, category: 'Electronics' },
  { id: '2', name: 'Mouse', price: 25, stock: 200, category: 'Electronics' },
  { id: '3', name: 'Keyboard', price: 75, stock: 100, category: 'Electronics' },
];

const validate = (values: EditProductFormData) => {
  const errors: Partial<EditProductFormData> = {};
  if (!values.name) {
    errors.name = 'Required';
  }
  if (!values.description) {
    errors.description = 'Required';
  }
  if (!values.price || values.price <= 0) {
    errors.price = 'Must be a positive number';
  }
  if (!values.stock || values.stock < 0) {
    errors.stock = 'Must be a non-negative number';
  }
  if (!values.category) {
    errors.category = 'Required';
  }
  return errors;
};

const onEditProductSubmit = async (values: EditProductFormData, router: any, productId: string) => {
  try {
    // TODO: Replace with actual API call to products microservice to update product
    console.log(`Updating product ${productId}:`, values);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000));

    alert(`Product "${values.name}" updated successfully!`);
    router.push('/seller/products');
  } catch (error: unknown) {
    alert(`Error updating product: ${(error as Error).message}`);
  }
};

const EditProductPage = ({ params }: { params: { id: string } }) => {
  const router = useRouter();
  const productId = params.id;
  const [initialValues, setInitialValues] = useState<EditProductFormData | null>(null);

  useEffect(() => {
    // In a real application, fetch product data from API using productId
    const product = mockProducts.find(p => p.id === productId);
    if (product) {
      setInitialValues(product);
    } else {
      // Handle product not found, e.g., redirect to 404 or products list
      alert('Product not found!');
      router.push('/seller/products');
    }
  }, [productId, router]);

  if (!initialValues) {
    return <div className="p-4">Loading product data...</div>;
  }

  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold mb-4">Edit Product: {initialValues.name}</h2>
      <Card>
        <Form
          onSubmit={(values) => onEditProductSubmit(values, router, productId)}
          validate={validate}
          initialValues={initialValues}
          render={({ handleSubmit, form, submitting, pristine }) => (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700">Product Name</label>
                <Field
                  name="name"
                  component="input"
                  type="text"
                  placeholder="Product Name"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                >
                  {({ input, meta }) => (
                    <div>
                      <input {...input} type="text" placeholder="Product Name" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                      {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                    </div>
                  )}
                </Field>
              </div>
              <div>
                <label htmlFor="description" className="block text-sm font-medium text-gray-700">Description</label>
                <Field
                  name="description"
                  component="textarea"
                  placeholder="Product Description"
                  rows={4}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                >
                  {({ input, meta }) => (
                    <div>
                      <textarea {...input} placeholder="Product Description" rows={4} className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                      {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                    </div>
                  )}
                </Field>
              </div>
              <div>
                <label htmlFor="price" className="block text-sm font-medium text-gray-700">Price</label>
                <Field
                  name="price"
                  component="input"
                  type="number"
                  placeholder="Price"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                >
                  {({ input, meta }) => (
                    <div>
                      <input {...input} type="number" placeholder="Price" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                      {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                    </div>
                  )}
                </Field>
              </div>
              <div>
                <label htmlFor="stock" className="block text-sm font-medium text-gray-700">Stock</label>
                <Field
                  name="stock"
                  component="input"
                  type="number"
                  placeholder="Stock"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                >
                  {({ input, meta }) => (
                    <div>
                      <input {...input} type="number" placeholder="Stock" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                      {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                    </div>
                  )}
                </Field>
              </div>
              <div>
                <label htmlFor="category" className="block text-sm font-medium text-gray-700">Category</label>
                <Field
                  name="category"
                  component="input"
                  type="text"
                  placeholder="Category"
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                >
                  {({ input, meta }) => (
                    <div>
                      <input {...input} type="text" placeholder="Category" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                      {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                    </div>
                  )}
                </Field>
              </div>
              <div className="flex space-x-4">
                <Button type="submit" disabled={submitting || pristine}>
                  Save Changes
                </Button>
                <Button type="button" onClick={() => router.push('/seller/products')} disabled={submitting}>
                  Cancel
                </Button>
              </div>
            </form>
          )}
        />
      </Card>
    </div>
  );
};

export default EditProductPage;
