'use client';

import { useEffect, useRef } from 'react';
import * as echarts from 'echarts';

export function TimelineChart() {
  const chartRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    if (!chartRef.current) return;
    
    const chart = echarts.init(chartRef.current);
    
    const statuses = ['Online', 'Idle', 'Do Not Disturb', 'Offline'];
    const colors = {
      'Online': '#3BA55C',
      'Idle': '#FAA61A',
      'Do Not Disturb': '#ED4245',
      'Offline': '#747F8D'
    };

    const timelineData = [
      [0, '2024-01-01 00:00:00', '2024-01-01 06:00:00', 'Offline'],
      [0, '2024-01-01 06:00:00', '2024-01-01 08:00:00', 'Idle'],
      [0, '2024-01-01 08:00:00', '2024-01-01 10:30:00', 'Online'],
      [0, '2024-01-01 10:30:00', '2024-01-01 11:00:00', 'Idle'],
      [0, '2024-01-01 11:00:00', '2024-01-01 12:30:00', 'Online'],
      [0, '2024-01-01 12:30:00', '2024-01-01 13:00:00', 'Do Not Disturb'],
      [0, '2024-01-01 13:00:00', '2024-01-01 15:30:00', 'Online'],
      [0, '2024-01-01 15:30:00', '2024-01-01 16:00:00', 'Do Not Disturb'],
      [0, '2024-01-01 16:00:00', '2024-01-01 16:30:00', 'Idle'],
      [0, '2024-01-01 16:30:00', '2024-01-01 17:30:00', 'Online'],
      [0, '2024-01-01 17:30:00', '2024-01-02 00:00:00', 'Offline']
    ];

    // Function to render the custom bars (timeline segments)
    function renderItem(params: any, api: any) {
      const categoryIndex = api.value(0);
      const start = api.coord([api.value(1), categoryIndex]);
      const end = api.coord([api.value(2), categoryIndex]);
      const height = api.size([0, 1])[1] * 0.6;

      const rectShape = echarts.graphic.clipRectByRect(
        {
          x: start[0],
          y: start[1] - height / 2,
          width: end[0] - start[0],
          height: height
        },
        {
          x: params.coordSys.x,
          y: params.coordSys.y,
          width: params.coordSys.width,
          height: params.coordSys.height
        }
      );

      return (
        rectShape && {
          type: 'rect',
          transition: ['shape'],
          shape: rectShape,
          style: {
            fill: colors[api.value(3) as keyof typeof colors]
          },
          styleEmphasis: {},
          textConfig: {
            position: 'inside'
          },
          textContent: {
            type: 'text',
            style: {
              text: api.value(3),
              fill: '#fff',
              fontSize: 12,
              fontWeight: 'bold',
              overflow: 'truncate',
            },
            invisible: end[0] - start[0] < 50
          }
        }
      );
    }

    const option = {
      tooltip: {
        formatter: function (params: any) {
          const startTime = echarts.format.formatTime('hh:mm', params.value[1]);
          const endTime = echarts.format.formatTime('hh:mm', params.value[2]);
          const status = params.value[3];
          const color = colors[status as keyof typeof colors];
          const customMarker = `<span style="display:inline-block;margin-right:5px;border-radius:50%;width:10px;height:10px;background-color:${color};"></span>`;
          return `${customMarker} ${status}: ${startTime} - ${endTime}`;
        }
      },
      legend: {
        data: statuses.map(status => ({
          name: status,
          itemStyle: {
            color: colors[status as keyof typeof colors]
          }
        })),
        selectedMode: false,
        bottom: 10
      },
      grid: {
        left: '5%',
        right: '5%',
        top: '10%',
        bottom: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'time',
        min: '2024-01-01 00:00:00',
        max: '2024-01-02 00:00:00',
        axisLabel: {
          formatter: function (value: any) {
            return echarts.format.formatTime('hh:00', value);
          }
        },
        interval: 3600 * 1000,
        axisTick: { show: true },
        axisLine: { show: false },
        splitLine: {
          show: true,
          lineStyle: {
            color: '#E5E7EB'
          }
        }
      },
      yAxis: {
        type: 'category',
        data: [],
        axisLine: { show: false },
        axisTick: { show: false }
      },
      series: [
        {
          type: 'custom',
          renderItem: renderItem,
          itemStyle: {
            opacity: 0.9,
            borderRadius: [5, 5, 5, 5]
          },
          label: {
            show: true,
            position: 'inside',
            formatter: '{@[3]}'
          },
          encode: {
            x: [1, 2],
            y: 0,
            tooltip: [1, 2, 3],
            itemName: 3
          },
          data: timelineData
        }
      ]
    };

    // Apply the configuration and render the chart
    chart.setOption(option);

    // Handle resize
    const handleResize = () => {
      chart.resize();
    };
    window.addEventListener('resize', handleResize);

    // Cleanup
    return () => {
      chart.dispose();
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  return <div ref={chartRef} className="w-full h-full" />;
} 