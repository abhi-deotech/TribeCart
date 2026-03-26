import React from 'react';
import Card from '@/components/ui/Card';

const SellerDashboardPage = () => {
  // Mock data for demonstration
  const totalProducts = 120;
  const totalOrders = 75;
  const totalSales = 15230.50;

  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold mb-4">Welcome to Your Seller Dashboard!</h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card title="Total Products">
          <p className="text-3xl font-semibold text-gray-800">{totalProducts}</p>
        </Card>
        <Card title="Total Orders">
          <p className="text-3xl font-semibold text-gray-800">{totalOrders}</p>
        </Card>
        <Card title="Total Sales">
          <p className="text-3xl font-semibold text-gray-800">${totalSales.toFixed(2)}</p>
        </Card>
      </div>

      <div className="mt-8">
        <h3 className="text-xl font-bold mb-3">Quick Actions</h3>
        {/* Add quick action buttons here, e.g., Add Product, View Orders */}
        <p>More dashboard content and quick actions will be added here.</p>
      </div>
    </div>
  );
};

export default SellerDashboardPage;