import { redirect } from 'next/navigation';

export default function Home() {
  redirect('/seller/login');
  return null;
}