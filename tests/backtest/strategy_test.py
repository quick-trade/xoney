# Copyright 2022 Vladyslav Kochetov. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# =============================================================================
import pandas as pd
from typing import Iterable

from xoney.strategy import Strategy
from xoney.generic.candlestick import Chart
from xoney.generic.events import Event, OpenTrade, CloseAllTrades
from xoney.generic.enums import TradeSide
from xoney.generic.trades import Trade
from xoney.generic.trades.levels import LevelHeap, SimpleEntry

from ta.volatility import BollingerBands


class BollingerTrendStrategy(Strategy):
    _signal = None
    _signal_prev_real = None

    def __init__(self, length=100, dev=1):
        self._settings["length"] = length
        self._settings["dev"] = dev

    def run(self, chart: Chart) -> None:
        bollinger = BollingerBands(
            close=pd.Series(chart.close[-self._settings["length"]:]),
            window=self._settings["length"],
            window_dev=self._settings["dev"]
        )
        signal_buy = bollinger.bollinger_hband().values[-1] < chart[-1].close
        signal_sell = bollinger.bollinger_lband().values[-1] > chart[-1].close

        signal = "long" if signal_buy else ("short" if signal_sell else None)

        if self._signal_prev_real != signal:
            self._signal = signal
        else:
            self._signal = None
        self._signal_prev_real = signal
        self.candle = chart[-1]

    def fetch_events(self) -> Iterable[Event]:
        if self._signal == "long":
            return [
                CloseAllTrades(),
                OpenTrade(
                    Trade(
                        side=TradeSide.LONG,
                        entries=LevelHeap(
                            [
                                SimpleEntry(
                                    trade_part=1,
                                    price=self.candle.close
                                )
                            ]
                        ),
                        breakouts=LevelHeap()
                    )
                )
            ]

        if self._signal == "short":
            return [
                CloseAllTrades(),
                OpenTrade(
                    Trade(
                        side=TradeSide.SHORT,
                        entries=LevelHeap(
                            [
                                SimpleEntry(
                                    trade_part=1,
                                    price=self.candle.close
                                )
                            ]
                        ),
                        breakouts=LevelHeap()
                    )
                )
            ]
        return []
