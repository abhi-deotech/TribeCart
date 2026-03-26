
import React from 'react';
import Link from 'next/link';

const Sidebar = () => {
  return (
    <div className="flex flex-col h-full bg-gray-800 text-white w-64">
      <div className="flex items-center justify-center h-16 border-b border-gray-700">
        <span className="text-xl font-semibold">Seller Dashboard</span>
      </div>
      <nav className="flex-1 px-2 py-4 space-y-2">
        <Link href="/seller" className="flex items-center px-3 py-2 rounded-md hover:bg-gray-700">
          Dashboard
        </Link>
        <Link href="/seller/products" className="flex items-center px-3 py-2 rounded-md hover:bg-gray-700">
          Products
        </Link>
        <Link href="/seller/orders" className="flex items-center px-3 py-2 rounded-md hover:bg-gray-700">
          Orders
        </Link>
        {/* Add more seller specific links here */}
      </nav>
    </div>
  );
};

export default Sidebar;
