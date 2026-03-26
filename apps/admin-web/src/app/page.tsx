"use client";

import { Button } from "@repo/ui/button";
import { useCounterStore } from "../store/counterStore";
import { useQuery } from "@tanstack/react-query";
import { Form, Field } from "react-final-form";

const fetchRandomNumber = async (): Promise<number> => {
  const response = await fetch("https://www.random.org/integers/?num=1&min=1&max=100&col=1&base=10&format=plain&rnd=new");
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  return parseInt(await response.text());
};

interface UserFormData {
  name: string;
  email: string;
  password: string;
}

const onRegisterSubmit = async (values: UserFormData) => {
  try {
    const response = await fetch("http://localhost:8000/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(values),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || "Registration failed");
    }

    const userData = await response.json();
    alert(`Registration successful! User ID: ${userData.id}`);
  } catch (error: unknown) {
    alert(`Error: ${(error as Error).message}`);
  }
};

export default function Home() {
  const { count, increment, decrement } = useCounterStore();
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["randomNumber"],
    queryFn: fetchRandomNumber,
  });

  if (isLoading) return <div>Loading random number...</div>;
  if (error) return <div>An error occurred: {error.message}</div>;

  return (
    <div className="flex flex-col items-center justify-center min-h-screen py-2">
      <h1 className="text-4xl font-bold mb-8">Admin Web</h1>
      <p className="text-2xl mb-4">Count: {count}</p>
      <div className="flex space-x-4 mb-8">
        <Button onClick={increment}>Increment</Button>
        <Button onClick={decrement}>Decrement</Button>
      </div>

      <p className="text-2xl mb-4">Random Number: {data}</p>
      <Button onClick={() => refetch()}>Fetch New Random Number</Button>

      <h2 className="text-3xl font-bold mt-8 mb-4">User Registration</h2>
      <Form
        onSubmit={onRegisterSubmit}
        initialValues={{ name: "", email: "", password: "" }}
        render={({ handleSubmit, form, submitting, pristine, values }) => (
          <form onSubmit={handleSubmit} className="flex flex-col space-y-4">
            <div>
              <label htmlFor="name" className="block text-lg font-medium text-gray-700">Name</label>
              <Field
                name="name"
                component="input"
                type="text"
                placeholder="Name"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
            <div>
              <label htmlFor="email" className="block text-lg font-medium text-gray-700">Email</label>
              <Field
                name="email"
                component="input"
                type="email"
                placeholder="Email"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
            <div>
              <label htmlFor="password" className="block text-lg font-medium text-gray-700">Password</label>
              <Field
                name="password"
                component="input"
                type="password"
                placeholder="Password"
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
            <div className="flex space-x-4">
              <Button type="submit" disabled={submitting || pristine}>
                Register
              </Button>
              <Button type="button" onClick={form.reset} disabled={submitting || pristine}>
                Reset
              </Button>
            </div>
            <pre className="mt-4 p-4 bg-gray-100 rounded-md text-sm">
              {JSON.stringify(values, undefined, 2)}
            </pre>
          </form>
        )}
      />
    </div>
  );
}
