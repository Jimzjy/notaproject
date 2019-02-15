import React from 'react'
import PropTypes from 'prop-types'
import { Color } from 'utils'
import CountUp from 'react-countup'
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import styles from './systemStats.less'

const countUpProps = {
  start: 0,
  duration: 2.75,
  useEasing: true,
  useGrouping: true,
  separator: ',',
}

function SystemStats({ data }) {
  return (
    <div className={styles.cpu}>
      <div className={styles.number}>
        <div className={styles.item}>
          <p>CPU Usage</p>
          <p>
            <CountUp end={data.cpuUsage} suffix="%" {...countUpProps} />
          </p>
        </div>
        <div className={styles.item}>
          <p>Memory Usage</p>
          <p>
            <CountUp end={data.memUsage} suffix="%" {...countUpProps} />
          </p>
        </div>
      </div>
      <ResponsiveContainer minHeight={360}>
        <AreaChart data={data.systemStats}>
          <Legend
            verticalAlign="top"
            align="right"
            iconType="circle"
            height={36}
          />
          <XAxis
            dataKey="name"
            axisLine={{ stroke: Color.borderBase, strokeWidth: 1 }}
            tickLine={false}
          />
          <YAxis axisLine={false} tickLine={false} />
          <CartesianGrid
            vertical={false}
            stroke={Color.borderBase}
            strokeDasharray="3 3"
          />
          <Tooltip
            wrapperStyle={{
              border: 'none',
              boxShadow: '4px 4px 40px rgba(0, 0, 0, 0.05)',
            }}
            content={content => {
              const list = content.payload.map((item, key) => (
                <li key={key} className={styles.tipitem}>
                  <span
                    className={styles.radiusdot}
                    style={{ background: item.color }}
                  />
                  {`${item.name}:${item.value}`}
                </li>
              ))
              return (
                <div className={styles.tooltip}>
                  <p className={styles.tiptitle}>{content.label}</p>
                  <ul>{list}</ul>
                </div>
              )
            }}
          />
          <Area
            name="Memory"
            type="monotone"
            dataKey="mem_used"
            stroke={Color.grass}
            fill={Color.grass}
            strokeWidth={2}
            dot={{ fill: '#fff' }}
            activeDot={{ r: 5, fill: '#fff', stroke: Color.green }}
          />
          <Area
            name="CPU"
            type="monotone"
            dataKey="cpu_used"
            stroke={Color.sky}
            fill={Color.sky}
            strokeWidth={2}
            dot={{ fill: '#fff' }}
            activeDot={{ r: 5, fill: '#fff', stroke: Color.blue }}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  )
}

SystemStats.propTypes = {
  data: PropTypes.object,
}

export default SystemStats
