"use client";

import React from 'react';
import { Form, Field } from 'react-final-form';
import { Button } from '@repo/ui/button';
import Card from '@/components/ui/Card';

interface SellerSignupFormData {
  email: string;
  password: string;
  confirmPassword?: string;
}

const validate = (values: SellerSignupFormData) => {
  const errors: Partial<SellerSignupFormData> = {};
  if (!values.email) {
    errors.email = 'Required';
  }
  if (!values.password) {
    errors.password = 'Required';
  }
  if (values.password !== values.confirmPassword) {
    errors.confirmPassword = 'Passwords must match';
  }
  return errors;
};

const onSellerSignupSubmit = async (values: SellerSignupFormData) => {
  try {
    // TODO: Replace with actual API call to users microservice for seller registration
    const response = await fetch("http://localhost:8000/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email: values.email, password: values.password }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || "Seller registration failed");
    }

    const userData = await response.json();
    alert(`Seller registration successful! User ID: ${userData.id}`);
    // TODO: Redirect to login or dashboard
  } catch (error: unknown) {
    alert(`Error: ${(error as Error).message}`);
  }
};

const SellerSignupPage = () => {
  return (
    <Card title="Seller Signup" className="w-full max-w-md p-8 space-y-6">
      <Form
        onSubmit={onSellerSignupSubmit}
        validate={validate}
        initialValues={{ email: '', password: '', confirmPassword: '' }}
        render={({ handleSubmit, form, submitting, pristine, values }) => (
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
            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700">Confirm Password</label>
              <Field
                name="confirmPassword"
                component="input"
                type="password"
                placeholder="Confirm Password"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                {({ input, meta }) => (
                  <div>
                    <input {...input} type="password" placeholder="Confirm Password" className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" />
                    {meta.error && meta.touched && <span className="text-red-500 text-xs mt-1">{meta.error}</span>}
                  </div>
                )}
              </Field>
            </div>
            <div className="flex space-x-4">
              <Button type="submit" disabled={submitting || pristine}>
                Sign Up
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

export default SellerSignupPage;
