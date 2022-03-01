import React from 'react';

const User = () => (
    <div>Users Module</div>
);

export default {
    routeProps: {
        path: '/user',
        component: User,
    },
    name: 'User',
};
