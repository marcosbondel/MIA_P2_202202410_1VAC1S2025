import { Navigate, Route, Routes } from 'react-router-dom'
import { LoginPage } from '../auth/pages'
import { MIAPage, TerminalPage } from '../mia/pages'
import { PublicRoute } from './PublicRoute'
import { PrivateRoute } from './PrivateRoute'
import { MIALayout } from '../mia/layout/MIALayout'
import { DiskPage } from '../mia/pages/DiskPage'
import { PartitionPage } from '../mia/pages/PartitionPage'

export const AppRouter = () => {
    return (
        <Routes>
            <Route path="/login" element={
                <PublicRoute>
                    <LoginPage />
                </PublicRoute>
            } />

            <Route path="/mia" element={
                <PrivateRoute>
                    <MIALayout />
                </PrivateRoute>
            }>
                <Route path="" element={<MIAPage />} />
                <Route path="terminal" element={<TerminalPage />} />
                <Route path="disks/:disk" element={<DiskPage />} />
                <Route path="disks/:disk/partitions/:partition" element={<PartitionPage />} />
                <Route path="*" element={<Navigate to='/mia' />} />
            </Route>

            <Route path="/*" element={<Navigate to='/mia' />}/>

        </Routes>
    )
}

