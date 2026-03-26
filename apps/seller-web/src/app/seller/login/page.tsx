"use client";

import React from 'react';
import { Form, Field } from 'react-final-form';
import { Button } from '@repo/ui/button';
import Card from '@/components/ui/Card';
import { useRouter } from 'next/navigation';

interface SellerLoginFormData {
  email: string;
  password: string;
}

const validate = (values: SellerLoginFormData) => {
  const errors: Partial<SellerLoginFormData> = {};
  if (!values.email) {
    errors.email = 'Required';
  }
  if (!values.password) {
    errors.password = 'Required';
  }
  return errors;
};

const onSellerLoginSubmit = async (values: SellerLoginFormData, router: any) => {
  try {
    // TODO: Replace with actual API call to users microservice for seller login
    const response = await fetch("http://localhost:8000/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: values.email, password: values.password }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || "Seller login failed");
    }

    const userData = await response.json();
    alert(`Seller login successful! Welcome, ${userData.email}`);
    // TODO: Store authentication token (e.g., in Zustand or local storage)
    router.push('/seller'); // Redirect to seller dashboard
  } catch (error: unknown) {
    alert(`Error: ${(error as Error).message}`);
  }
};

const SellerLoginPage = () => {
  const router = useRouter();

  return (
    <Card title="Seller Login" className="w-full max-w-md p-8 space-y-6">
      <Form
        onSubmit={(values) => onSellerLoginSubmit(values, router)}
        validate={validate}
        initialValues={{ email: '', password: '' }}
        render={({ handleSubmit, form, submitting, pristine }) => (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">Email</label>
              <Field
                name="email"
                component="input"
                type="email"
                placeholder="Email"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                {({ input, meta }) => (
                  <div>
                    <input {...input} type="email" placeholder="Email" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                    {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                  </div>
                )}
              </Field>
            </div>
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700">Password</label>
              <Field
                name="password"
                component="input"
                type="password"
                placeholder="Password"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                {({ input, meta }) => (
                  <div>
                    <input {...input} type="password" placeholder="Password" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                    {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                  </div>
                )}
              </Field>
            </div>
            <div className="flex space-x-4">
              <Button type="submit" disabled={submitting || pristine}>
                Login
              </Button>
              <Button type="button" onClick={form.reset} disabled={submitting || pristine}>
                Reset
              </Button>
            </div>
          </form>
        )}
      />
    </Card>
  );
};

export default SellerLoginPage;
