# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class BzzChunksPinnedChunks(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, address=None, pin_counter=None):  # noqa: E501
        """BzzChunksPinnedChunks - a model defined in OpenAPI

        :param address: The address of this BzzChunksPinnedChunks.  # noqa: E501
        :type address: str
        :param pin_counter: The pin_counter of this BzzChunksPinnedChunks.  # noqa: E501
        :type pin_counter: int
        """
        self.openapi_types = {
            'address': str,
            'pin_counter': int
        }

        self.attribute_map = {
            'address': 'address',
            'pin_counter': 'pinCounter'
        }

        self._address = address
        self._pin_counter = pin_counter

    @classmethod
    def from_dict(cls, dikt) -> 'BzzChunksPinnedChunks':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The BzzChunksPinned_chunks of this BzzChunksPinnedChunks.  # noqa: E501
        :rtype: BzzChunksPinnedChunks
        """
        return util.deserialize_model(dikt, cls)

    @property
    def address(self):
        """Gets the address of this BzzChunksPinnedChunks.


        :return: The address of this BzzChunksPinnedChunks.
        :rtype: str
        """
        return self._address

    @address.setter
    def address(self, address):
        """Sets the address of this BzzChunksPinnedChunks.


        :param address: The address of this BzzChunksPinnedChunks.
        :type address: str
        """

        self._address = address

    @property
    def pin_counter(self):
        """Gets the pin_counter of this BzzChunksPinnedChunks.


        :return: The pin_counter of this BzzChunksPinnedChunks.
        :rtype: int
        """
        return self._pin_counter

    @pin_counter.setter
    def pin_counter(self, pin_counter):
        """Sets the pin_counter of this BzzChunksPinnedChunks.


        :param pin_counter: The pin_counter of this BzzChunksPinnedChunks.
        :type pin_counter: int
        """

        self._pin_counter = pin_counter
